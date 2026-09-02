//go:build ignore

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("WAILS_VITE_PORT")
	if port == "" {
		port = "9245"
	}
	fs := http.FileServer(http.Dir("frontend/dist"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Prevent caching in dev mode so edits reflect immediately
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fs.ServeHTTP(w, r)
	})

	fmt.Printf("[DevServer] Serving frontend/dist on http://localhost:%s\n", port)
	if err := http.ListenAndServe("127.0.0.1:"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "[DevServer] Error: %v\n", err)
		os.Exit(1)
	}
}
