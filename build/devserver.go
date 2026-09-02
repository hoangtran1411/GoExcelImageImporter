//go:build ignore

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := os.Getenv("WAILS_VITE_PORT")
	if port == "" {
		port = "9245"
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("frontend/dist"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Prevent caching in dev mode so edits reflect immediately
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fs.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    "127.0.0.1:" + port,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("[DevServer] Serving frontend/dist on http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "[DevServer] Error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-stop
	fmt.Println("\n[DevServer] Shutting down cleanly...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "[DevServer] Shutdown error: %v\n", err)
	}
}
