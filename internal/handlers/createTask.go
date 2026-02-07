package handlers

import (
	"net/http"
	"time"
	"strings"
)

func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.CreateTask"

	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", "fn", fn, "err", err)
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		h.logger.Error("Missing required field 'name'", "fn", fn)
    http.Error(w, "Field 'name' is required", http.StatusBadRequest)
    return
	}

	query := `INSERT INTO Tasks (Name, Completed, Created_at, Completed_at)
						VALUES (?, ?, ?, ?) `

	res, err := h.storage.DB.Exec(query, name, false, time.Now(), nil)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			h.logger.Error("Task with this name already exists", "fn", fn, "task_name", name)
			http.Error(w, "Task with this name already exists", http.StatusBadRequest)
			return
		}

		h.logger.Error("Insert query error", "fn", fn, "query", query, "err", err)
		http.Error(w, "Insert query error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		h.logger.Error("Failed to get last insert id", "fn", fn, "err", err)
		return
	}

	h.logger.Info("Task created successfully", "fn", fn, "task_id", id, "task_name", name)
	
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
	h.logger.Info("Redirect success", "fn", fn)
}