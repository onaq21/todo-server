package main

import (
	"net/http"
	"os"
	"html/template"
	"log"
	"time"
	"github.com/onaq21/todo-server/internal/config"
	"github.com/onaq21/todo-server/internal/logger"
	"github.com/onaq21/todo-server/internal/storage/sqlite"
	"github.com/onaq21/todo-server/internal/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %s", err)
	}
	log.Println("config load success")

	logg := logger.New(cfg.Env)
	logg.Info("logger initialized", "env", cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logg.Error("open db error", "err", err)
		os.Exit(1)
	}
	defer storage.DB.Close()
	logg.Info("open db success")

	funcs := template.FuncMap(map[string]any{
		"formatTime": func(t time.Time) string {
			if t.IsZero() { return "-" }
			return t.Format("02.01.2006 15:04")
		},
	})
	logg.Info("template functions registered", "functions", []string{"formatTime"})

	tmplTasks := template.Must(template.New("tasks").Funcs(funcs).ParseFiles("internal/templates/base.html", "internal/templates/tasks.html"))
	logg.Info("Tasks template parsed successfully", "files", []string{"base.html", "tasks.html"})

	tmplEdit := template.Must(template.New("edit").Funcs(funcs).ParseFiles("internal/templates/base.html", "internal/templates/edit.html"))
	logg.Info("Edit template parsed successfully", "files", []string{"base.html", "edit.html"})

	h := handlers.NewHandler(storage, logg, tmplTasks, tmplEdit)

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/", http.RedirectHandler("/tasks", http.StatusFound))
	mux.HandleFunc("GET /tasks", h.GetAllTasksHandler)
	mux.HandleFunc("GET /tasks/{id}/edit", h.GetTaskHandler)
	mux.HandleFunc("GET /tasks/sort", h.SortTasksHandler)
	mux.HandleFunc("POST /tasks", h.CreateTaskHandler)
	mux.HandleFunc("POST /tasks/{id}/edit", h.UpdateTaskHandler)
	mux.HandleFunc("POST /tasks/{id}/delete", h.DeleteTaskHandler)

	logg.Info("starting HTTP server", "address", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, mux); err != nil {
		logg.Error("HTTP server failed to start", "address", cfg.Address, "err", err)
		os.Exit(1)
	}
}