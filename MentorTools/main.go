package main

import (
	"context"
	"log"
	"net/http"

	"MentorTools/db"
	"MentorTools/handlers"
	"MentorTools/middleware"

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

	// Защищённые маршруты - только для авторизованных пользователей
	router.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler(conn)))).Methods("GET")

	// Добавляем маршрут для поиска пользователей в зависимости от роли
	router.Handle("/search", middleware.AuthMiddleware(http.HandlerFunc(handlers.SearchUsersHandler(conn)))).Methods("GET")

	// Добавляем маршрут для сохранения измененых данных о пользователе
	router.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateProfileHandler(conn)))).Methods("PUT")

	// Роут для создания связи между учителем и учеником
	router.Handle("/link", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateTeacherStudentLink(conn)))).Methods("POST")

	// Роут для удаления связи
	router.Handle("/unlink", middleware.AuthMiddleware(http.HandlerFunc(handlers.RemoveTeacherStudentLink(conn)))).Methods("DELETE")

	// Роут для отображения всех связей
	router.Handle("/links", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetTeacherStudentLinks(conn)))).Methods("GET")

	// Статические файлы
	fs := http.FileServer(http.Dir("./fe"))
	router.PathPrefix("/").Handler(fs)

	// Запуск сервера
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
