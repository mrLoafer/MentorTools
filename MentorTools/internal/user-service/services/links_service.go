package services

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/pkg/common"
	"context"
)

// Функция для получения информации об учителе
func GetTeacherInfo(ctx context.Context, claims jwt.MapClaims) (models.User, *common.Response) {
	studentID := int(claims["userId"].(float64))

	teacherInfo, err := getTeacherFromDB(ctx, studentID)
	if err != nil {
		return models.User{}, common.NewErrorResponse("TEACHER_NOT_FOUND", "Teacher not found")
	}

	return teacherInfo, nil
}

// Функция для получения списка учеников
func GetStudents(ctx context.Context, claims jwt.MapClaims) ([]models.User, *common.Response) {
	teacherID := int(claims["userId"].(float64))

	students, err := getStudentsFromDB(ctx, teacherID)
	if err != nil {
		return nil, common.NewErrorResponse("STUDENTS_NOT_FOUND", "No students found")
	}

	return students, nil
}

// Функция для создания связи с учеником по email
func CreateLink(ctx context.Context, claims jwt.MapClaims, studentEmail string) *common.Response {
	teacherID := int(claims["userId"].(float64))

	// Валидация роли учителя и существования студента
	studentID, err := getStudentIDByEmail(ctx, studentEmail)
	if err != nil {
		return common.NewErrorResponse("STUDENT_NOT_FOUND", "Student not found")
	}

	if err := createLinkInDB(ctx, teacherID, studentID); err != nil {
		return common.NewErrorResponse("LINK_CREATION_FAILED", "Failed to create link")
	}

	return nil
}
