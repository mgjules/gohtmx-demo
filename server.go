package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mgjules/gohtmx-demo/task"
	"github.com/mgjules/gohtmx-demo/templates"
	uuid "github.com/satori/go.uuid"
)

func newServer(addr string, manager *task.Manager) (*http.Server, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be empty")
	}
	if manager == nil {
		return nil, errors.New("manager cannot be nil")
	}

	mux := http.NewServeMux()

	// Routes.
	mux.HandleFunc("GET /", handleIndex(manager))
	mux.HandleFunc("GET /tasks", handleListTask(manager))
	mux.HandleFunc("POST /tasks", handleAddTask(manager))
	mux.HandleFunc("DELETE /tasks/{id}/done", handleMarkTaskAsDone(manager))
	mux.Handle("GET /assets/dist/", http.FileServerFS(assets))

	return &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 6,
	}, nil
}

func handleIndex(manager *task.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := templates.IndexPage(manager.List()).Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render index page", http.StatusInternalServerError)
		}
	}
}

func handleAddTask(manager *task.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			errMsg  string
			content = r.FormValue("task")
			task    *task.Task
		)
		defer func() {
			w.Header().Add("Content-Type", "text/html")
			if errMsg == "" && task != nil {
				if err := templates.TaskWrappedItemComponent(*task).Render(r.Context(), w); err != nil {
					slog.Error("Failed to render task", "error", err)
					errMsg = "Failed to render task." + err.Error() + "."
				}
			}
			if err := templates.TaskInputComponent(content, errMsg).Render(r.Context(), w); err != nil {
				slog.Error("Failed to add task", "error", err)
				http.Error(w, "Failed to render task input component", http.StatusInternalServerError)
			}
		}()

		if content == "" {
			slog.Error("Failed to retrieve task content: empty")
			errMsg = "Please input a task"
			return
		}

		if task, err = manager.Add(content); err != nil {
			slog.Error("Failed to add task", "error", err)
			errMsg = err.Error() + "."
			return
		}
	}
}

func handleListTask(manager *task.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		if err := templates.TaskListComponent(manager.List()).Render(r.Context(), w); err != nil {
			http.Error(w, "Failed to render task list component", http.StatusInternalServerError)
		}
	}
}

func handleMarkTaskAsDone(manager *task.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := r.PathValue("id")
		id, err := uuid.FromString(raw)
		if err != nil {
			slog.Error("Failed to parse id", "id", raw, "error", err)
			http.Error(w, fmt.Sprintf("Invalid ID %q", raw), http.StatusBadRequest)
			return
		}

		if err := manager.MarkAsDone(id); err != nil {
			slog.Error("Failed to mark task as done", "id", id, "error", err)
			http.Error(w, fmt.Sprintf("Failed to mark task %q as done", id.String()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
