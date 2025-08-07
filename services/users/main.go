package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service":"users","status":"ok"}`)
	})

	log.Println("Users service starting on :8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
}
