package main

import (
	"log/slog"
	"os"
)

func initLogger(prod bool) {
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if prod {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
