package handlers

import (
	"github.com/onaq21/todo-server/internal/storage/sqlite"
	"log/slog"
	"html/template"
)

type Handler struct {
	storage *sqlite.Storage
	logger *slog.Logger
	tmplTasks *template.Template
	tmplEdit *template.Template 
}

func NewHandler(storage *sqlite.Storage, logger *slog.Logger, tmplTasks *template.Template, tmplEdit *template.Template) *Handler {
	return &Handler{storage, logger, tmplTasks, tmplEdit}
}