package handlers

import (
	"net/http"
	"database/sql"
	"github.com/onaq21/todo-server/internal/task"
)

func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.GetTask"

	query := `SELECT * FROM Tasks WHERE ID = ?`
	id := r.PathValue("id")

	row := h.storage.DB.QueryRow(query, id)

	var task task.Task

	if err := row.Scan(&task.ID, &task.Name, &task.Completed, &task.CreatedAt, &task.CompletedAt); err != nil {
		if err == sql.ErrNoRows {
			h.logger.Error("Task not found", "fn", fn, "task_id", id)
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			h.logger.Error("Query error", "fn", fn, "query", query, "task_id", id, "err", err)
			http.Error(w, "Query error: " + err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmplEdit.ExecuteTemplate(w, "base.html", task); err != nil {
		h.logger.Error("template execution error", "fn", fn, "err", err)
		http.Error(w, "Template execution error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("GetTaskHandler executed successfully", "fn", fn, "task_id", task.ID)
}
