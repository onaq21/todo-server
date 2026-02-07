package handlers

import (
	"net/http"
	"github.com/onaq21/todo-server/internal/task"
)

func (h *Handler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.GetAllTasks"
	
	query := `SELECT * FROM Tasks ORDER BY ID`

	rows, err := h.storage.DB.Query(query)
	if err != nil {
		h.logger.Error("Select db error", "fn", fn, "query", query, "err", err)
		http.Error(w, "Select db error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := make([]task.Task, 0)

	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed, &task.CreatedAt, &task.CompletedAt); err != nil {
			h.logger.Error("Scan row error", "fn", fn, "err", err)
			http.Error(w, "Scan error: " + err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		h.logger.Error("Rows iteration error", "fn", fn, "err", err)
		http.Error(w, "Error iterating over rows: " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmplTasks.ExecuteTemplate(w, "base.html", tasks); err != nil {
		h.logger.Error("template execution error", "fn", fn, "err", err)
		http.Error(w, "Template execution error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("GetAllTasksHandler executed successfully", "fn", fn, "tasks_count", len(tasks))
}
