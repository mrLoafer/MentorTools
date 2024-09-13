package main

import (
	"context"
	"log"
	"net/http"

	"MentorTools/db"
	"MentorTools/handlers"
	"MentorTools/middleware"
	"MentorTools/users"

	"github.com/gorilla/mux"
)

func main() {
	// Подключение к базе данных
	conn := db.ConnectDB()
	defer conn.Close(context.Background())

	// Маршруты
	http.HandleFunc("/login", handlers.LoginHandler(conn))
	http.HandleFunc("/register", handlers.RegisterHandler(conn))
	http.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProtectedHandler)))

	router := mux.NewRouter()

	// Маршруты для управления пользователями
	router.HandleFunc("/users/{id}", users.GetUserHandler(conn)).Methods("GET")
	router.HandleFunc("/users/{id}", users.UpdateUserHandler(conn)).Methods("PUT")
	router.HandleFunc("/users", users.ListUsersHandler(conn)).Methods("GET")

	// Запуск сервера
	// Раздача статических файлов (HTML, CSS, JS) из папки "static"
	fs := http.FileServer(http.Dir("./fe"))
	http.Handle("/", fs) // Все запросы по корневому URL направляются в папку static

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
