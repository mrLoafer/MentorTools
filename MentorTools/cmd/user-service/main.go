package main

import (
	"MentorTools/internal/user-service/handlers"
	"MentorTools/internal/user-service/repository"
	"MentorTools/pkg/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	// Initialize database connection
	dbPool, err := repository.InitDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repository.CloseDB(dbPool)

	// Setting up routes with middleware for authorization
	http.Handle("/dashboard", middleware.AuthMiddleware(handlers.DashboardHandler()))
	http.Handle("/profile", middleware.AuthMiddleware(handlers.UpdateUserProfileHandler(dbPool)))
	http.Handle("/students", middleware.AuthMiddleware(handlers.GetStudentsHandler(dbPool)))
	http.Handle("/link", middleware.AuthMiddleware(handlers.CreateLinkHandler(dbPool)))

	// Health check route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User-service is running")
	})

	// Start the server
	fmt.Println("Starting user-service on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
