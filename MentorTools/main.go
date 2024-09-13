package main

import (
	"context"
	"log"
	"net/http"

	"MentorTools/db"
	"MentorTools/handlers"
	"MentorTools/users"

	"github.com/gorilla/mux"
)

func main() {
	// Подключение к базе данных
	conn := db.ConnectDB()
	defer conn.Close(context.Background())

	// Создание роутера
	router := mux.NewRouter()

	// Маршруты для авторизации и регистрации
	router.HandleFunc("/login", handlers.LoginHandler(conn)).Methods("POST")
	router.HandleFunc("/register", handlers.RegisterHandler(conn)).Methods("POST")

	// Маршруты для управления пользователями
	router.HandleFunc("/users/{id}", users.GetUserHandler(conn)).Methods("GET")
	router.HandleFunc("/users/{id}", users.UpdateUserHandler(conn)).Methods("PUT")
	router.HandleFunc("/users", users.ListUsersHandler(conn)).Methods("GET")

	// Раздача статических файлов (HTML, CSS, JS) из папки "fe"
	fs := http.FileServer(http.Dir("./fe"))
	router.PathPrefix("/").Handler(fs) // Все запросы по корневому URL направляются в папку fe

	// Запуск сервера
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
