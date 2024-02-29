package main

import (
	"log/slog"
	"os"
)

func initLogger(prod bool) {
	var handler slog.Handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	if prod {
		handler = slog.NewJSONHandler(os.Stderr, nil)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
