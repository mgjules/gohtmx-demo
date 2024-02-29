package main

import (
	"net/http"
	"time"
)

func newServer(addr string) *http.Server {
	mux := http.NewServeMux()

	// Routes.
	mux.HandleFunc("GET /", handleIndex)
	mux.Handle("GET /assets/dist/", http.FileServerFS(assets))

	return &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 6,
	}
}
