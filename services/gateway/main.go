package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service":"gateway","status":"ok"}`)
	})

	log.Println("Gateway starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
