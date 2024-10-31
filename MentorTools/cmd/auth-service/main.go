package main

import (
	"MentorTools/internal/auth-service/handlers"
	"MentorTools/internal/auth-service/repository"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	// Initialize the database connection
	if err := repository.InitDB(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repository.CloseDB()

	// Register and auth routes
	http.HandleFunc("/register", handlers.RegisterHandler())
	http.HandleFunc("/login", handlers.LoginHandler())

	// Health check route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Service is running")
	})

	// Start the server
	fmt.Println("Starting auth-service on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
