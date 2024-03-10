// GoHTMX is just a simple application showcasing Go + HTMX and some Templ.
package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mgjules/gohtmx-demo/task"
)

//go:embed assets/dist
var assets embed.FS

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to start service", "error", err)
		os.Exit(1)
	}
}

func run() error {
	prod := flag.Bool("prod", false, "are we running in production")
	addr := flag.String("addr", "localhost:8080", "address of http server")
	flag.Parse()

	initLogger(*prod)

	gofakeit.Seed(13337)

	manager := task.NewManager()
	if err := manager.Seed(5); err != nil {
		return fmt.Errorf("failed to seed tasks: %v", err)
	}

	server, err := newServer(*addr, manager)
	if err != nil {
		return fmt.Errorf("failed to create new server: %w", err)
	}

	slog.Info("Server starting...", "Address", "http://"+server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen and serve http server: %w", err)
	}

	return nil
}
