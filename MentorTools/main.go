package main

import (
	"context"
	"log"
	"net/http"

	"MentorTools/handlers"
	"MentorTools/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Подключаемся к базе данных через пул соединений
	dbpool, err := pgxpool.Connect(context.Background(), "postgresql://loafer:Tesla846@localhost:5432/mentor_tools")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Создание роутера
	router := mux.NewRouter()

	// Маршруты для авторизации и регистрации
	router.HandleFunc("/login", handlers.LoginHandler(dbpool)).Methods("POST")
	router.HandleFunc("/register", handlers.RegisterHandler(dbpool)).Methods("POST")

	// Защищённые маршруты - только для авторизованных пользователей
	router.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfileHandler(dbpool)))).Methods("GET")

	// Добавляем маршрут для поиска пользователей в зависимости от роли
	router.Handle("/search", middleware.AuthMiddleware(http.HandlerFunc(handlers.SearchUsersHandler(dbpool)))).Methods("GET")

	// Добавляем маршрут для сохранения измененых данных о пользователе
	router.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateProfileHandler(dbpool)))).Methods("PUT")

	// Роут для топиков
	router.Handle("/get-contexts", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetAvailableContextsHandler(dbpool)))).Methods("GET")

	// Роут для создания связи между учителем и учеником
	router.Handle("/link", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateTeacherStudentLink(dbpool)))).Methods("POST")

	// Роут для удаления связи
	router.Handle("/unlink", middleware.AuthMiddleware(http.HandlerFunc(handlers.RemoveTeacherStudentLink(dbpool)))).Methods("DELETE")

	// Роут для отображения всех связей
	router.Handle("/links", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetTeacherStudentLinks(dbpool)))).Methods("GET")

	//Роуты для работы с словарем
	router.Handle("/get-words", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetWordsHandler(dbpool)))).Methods("GET")
	router.Handle("/add-word", middleware.AuthMiddleware(http.HandlerFunc(handlers.AddWordHandler(dbpool)))).Methods("POST")
	router.Handle("/update-word-status", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateWordStatusHandler(dbpool)))).Methods("POST")
	router.Handle("/get-word-details", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetWordDetailsHandler(dbpool)))).Methods("GET")

	// Маршрут для получения роли пользователя
	router.Handle("/get-user-role", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserRoleHandler()))).Methods("GET")

	// Статические файлы
	fs := http.FileServer(http.Dir("./fe"))
	router.PathPrefix("/").Handler(fs)

	// Запуск сервера
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
