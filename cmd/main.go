package main

import (
	"log"
	"net/http"
	
	handlers "movieshop/backend/Handlers"
)

func main() {
	// Create handler with routes
	handler := handlers.MovieListHandler()
	
	// Start server
	log.Println("Starting server on :8000")
	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
