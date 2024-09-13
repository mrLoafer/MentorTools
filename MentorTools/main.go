package main

import (
	"context"
	"log"
	"mentor_tools/db"
	"mentor_tools/handlers"
	"mentor_tools/middleware"
	"net/http"
)

func main() {
	// Подключение к базе данных
	conn := db.ConnectDB()
	defer conn.Close(context.Background())

	// Маршруты
	http.HandleFunc("/login", handlers.LoginHandler(conn))
	http.HandleFunc("/register", handlers.RegisterHandler(conn))
	http.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProtectedHandler)))

	// Запуск сервера
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
