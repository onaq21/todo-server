package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/onaq21/todo-server/internal/task"
)

type AIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type AISortResponse struct {
	OrderedIDs []int `json:"ordered_ids"`
}

func SortTasksByMeaning(tasks []task.Task) ([]task.Task, error) {
	const fn = "internal.ai.Sort"

	prompt := `You will receive a list of tasks.
	Each task has:
	- id: integer (unique, must be preserved exactly)
	- name: string
	- completed: boolean
	SORTING RULES (STRICT):
	1. completed = false tasks go first
	2. completed = true tasks go last
	3. DO NOT remove tasks
	4. DO NOT change IDs
	5. KEEP ALL TASKS
	OUTPUT (STRICT JSON ONLY):
	{
		"ordered_ids": [1, 5, 3]
	}
	INPUT TASKS:
	`

	for _, t := range tasks {
		prompt += fmt.Sprintf("- id: %d, name: \"%s\", completed: %t\n", t.ID, t.Name, t.Completed,)
	}

	reqBody := AIRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []Message{
			{
				Role: "system",
				Content: "You are an assistant that sorts tasks by priority. You MUST respond ONLY with valid JSON in the format {\"ordered_ids\": [1,2,3]} and nothing else.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   500,
		Temperature: 0.0,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to marshal request: %w", fn, err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create HTTP request: %w", fn, err)
	}

	token := os.Getenv("GROQ_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("%s: GROQ_API_KEY is not set", fn)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", fn, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody bytes.Buffer
		errBody.ReadFrom(resp.Body)
		return nil, fmt.Errorf("%s: AI returned status %d: %s", fn, resp.StatusCode, errBody.String())
	}

	var aiResp AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, fmt.Errorf("%s: failed to decode AI response: %w", fn, err)
	}

	if len(aiResp.Choices) == 0 {
		return nil, fmt.Errorf("%s: empty AI response", fn)
	}

	rawContent := aiResp.Choices[0].Message.Content
	var sortedNames AISortResponse
	if err := json.Unmarshal([]byte(rawContent), &sortedNames); err != nil {
		return nil, fmt.Errorf("%s: failed to unmarshal AI content: %w; raw: %s", fn, err, rawContent)
	}

	taskMap := make(map[int]task.Task, len(tasks))
	for _, t := range tasks {
		taskMap[t.ID] = t
	}

	sortedTasks := make([]task.Task, 0, len(tasks))
	seen := make(map[int]bool, len(tasks))

	for _, id := range sortedNames.OrderedIDs {
		if t, ok := taskMap[id]; ok {
			sortedTasks = append(sortedTasks, t)
			seen[id] = true
		}
	}

	for _, t := range tasks {
		if !seen[t.ID] {
			sortedTasks = append(sortedTasks, t)
		}
	}

	return sortedTasks, nil
}