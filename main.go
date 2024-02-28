// GoHTMX is just a simple application showcasing Go + HTMX and some Templ.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/mgjules/gohtmx-demo/templates"
)

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

	server := newServer(*addr)

	slog.Info("Server starting...", "Address", "http://"+server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen and serve http server: %w", err)
	}

	return nil
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := templates.Hello("testing").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render Hello template", http.StatusInternalServerError)
	}
}
