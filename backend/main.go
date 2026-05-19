package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go Backend!")
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	port := "8080"
	fmt.Printf("Backend server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
