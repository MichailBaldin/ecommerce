package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service":"products","status":"ok"}`)
	})

	log.Println("Products service starting on :8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
