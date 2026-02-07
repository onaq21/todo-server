package handlers

import (
	"net/http"
	"strconv"
	"time"
	"strings"
)

func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.UpdateTask"

	id := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", "fn", fn, "err", err)
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	val := r.FormValue("completed")
	completed, err := strconv.ParseBool(val)
	if err != nil {
		h.logger.Error("Invalid completed value", "fn", fn, "val", val, "err", err)
		http.Error(w, "Invalid completed value", http.StatusBadRequest)
		return
	}

	query := "UPDATE Tasks SET "
	sets := []string{"Name = ?", "Completed = ?"}
	args := []any{name, completed}

	if completed {
			args = append(args, time.Now())
			sets = append(sets, "Completed_at = ?")
	} else {
			args = append(args, nil)
			sets = append(sets, "Completed_at = ?")
	}

	query += strings.Join(sets, ", ")
	query += " WHERE ID = ?"
	args = append(args, id)

	res, err := h.storage.DB.Exec(query, args...)
	if err != nil {
		h.logger.Error("Update error", "fn", fn, "query", query, "err", err)
		http.Error(w, "Update error", http.StatusInternalServerError)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		h.logger.Error("RowsAffected error", "fn", fn, "task_id", id, "err", err)
		http.Error(w, "RowsAffected error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	h.logger.Info("Task updated successfully", "fn", fn, "task_id", id)

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	h.logger.Info("Redirect success", "fn", fn)
}