package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidsugianto/go-pkgs/grace"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! Time: %s\n", time.Now().Format(time.RFC3339))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})

	fmt.Println("Starting server on :8080")
	fmt.Println("Try: curl http://localhost:8080/slow")
	fmt.Println("Then press Ctrl+C to see graceful shutdown")

	grace.ServeHTTP(":8080", nil)
}
