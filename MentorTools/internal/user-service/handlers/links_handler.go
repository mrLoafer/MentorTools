package handlers

import (
	"MentorTools/internal/user-service/models"
	"MentorTools/internal/user-service/services"
	"MentorTools/pkg/common"
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
)

// GetTeacherInfoHandler обрабатывает запрос на получение информации об учителе
func GetTeacherInfoHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok || claims["role"] != "student" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(common.NewErrorResponse("AUTH403", "Access denied"))
		return
	}

	teacherInfo, appErr := services.GetTeacherInfo(context.Background(), claims)
	if appErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.NewErrorResponse("USER404", appErr.Message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.NewSuccessResponse("Teacher information", teacherInfo))
}

// GetStudentsHandler обрабатывает запрос на получение списка учеников (доступно только для учителей)
func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok || claims["role"] != "teacher" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(common.NewErrorResponse("AUTH403", "Access denied"))
		return
	}

	students, appErr := services.GetStudents(context.Background(), claims)
	if appErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(common.NewErrorResponse("USER404", appErr.Message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.NewSuccessResponse("Students list", students))
}

// CreateLinkHandler обрабатывает запрос на установку связи между учителем и учеником
func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(jwt.MapClaims)
	if !ok || claims["role"] != "teacher" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(common.NewErrorResponse("AUTH403", "Access denied"))
		return
	}

	var linkRequest models.UserLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&linkRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.NewErrorResponse("REQ400", "Invalid request payload"))
		return
	}

	appErr := services.CreateLink(context.Background(), claims, linkRequest.Email)
	if appErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(common.NewErrorResponse("LINK409", appErr.Message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.NewSuccessResponse("Link created successfully", nil))
}
