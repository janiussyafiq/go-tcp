package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Simple backend server for testing
func startBackend(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf("Backend server on port %s received request to %s", port, r.URL.Path)

		// Print all headers received
		response += "\nHeaders received:\n"
		for name, values := range r.Header {
			for _, value := range values {
				response += fmt.Sprintf(" %s: %s\n", name, value)
			}
		}

		fmt.Println(response)
		w.Write([]byte(response))
	})

	log.Printf("Backend server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run backend.go <port>")
		os.Exit(1)
	}

	port := os.Args[1]
	startBackend(port)
}
