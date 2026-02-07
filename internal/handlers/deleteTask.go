package handlers

import (
	"net/http"
)

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.deleteTask"
	
	query := `DELETE FROM Tasks WHERE ID = ?`
	id := r.PathValue("id")

	res, err := h.storage.DB.Exec(query, id)
	if err != nil {
		h.logger.Error("Query error", "fn", fn, "query", query, "task_id", id, "err", err)
		http.Error(w, "Delete query error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		h.logger.Error("RowsAffected error", "fn", fn, "task_id", id, "err", err)
		http.Error(w, "RowsAffected error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		h.logger.Error("Task not found", "fn", fn, "task_id", id)
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	h.logger.Info("Task deleted successfully", "fn", fn, "task_id", id)
	
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	h.logger.Info("Redirect success", "fn", fn)
}