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

	// Initialize the database connection and store it in a local variable
	dbPool, err := repository.InitDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repository.CloseDB(dbPool)

	// Register and auth routes with injected dbPool
	http.HandleFunc("/register", handlers.RegisterHandler(dbPool))
	http.HandleFunc("/login", handlers.LoginHandler(dbPool))

	// Health check route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Auth-service is running")
	})

	// Start the server
	fmt.Println("Starting auth-service on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
