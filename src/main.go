package main

import (
	// Standard library imports
	"fmt"
	"log"
	"net/http"
	"os"

	// Internal Imports
	"log-forwarder/handler"
)

func main() {
	fmt.Println("Starting server...")

	http.HandleFunc("/messages", handler.PubSubHandler)
	http.HandleFunc("/hello", handler.HealthCheckHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Printf("Defaulting to port %s\n", port)
	}

	fmt.Printf("Log Forwarder listening on port: '%s'\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
