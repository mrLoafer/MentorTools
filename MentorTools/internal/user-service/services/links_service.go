package services

import (
	"MentorTools/internal/user-service/models"
	"MentorTools/pkg/common"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// GetStudents retrieves the list of students linked to a specific teacher.
func GetStudents(ctx context.Context, dbPool *pgxpool.Pool, teacherID int) ([]models.Student, *common.Response) {
	var students []models.Student

	// Call the database procedure to get students linked to the teacher
	rows, err := dbPool.Query(ctx, "SELECT * FROM fn_get_students_by_teacher($1)", teacherID)
	if err != nil {
		return nil, common.NewErrorResponse("DB500", "Database error occurred while retrieving students")
	}
	defer rows.Close()

	// Scan each row into the Student model
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email); err != nil {
			return nil, common.NewErrorResponse("DB500", "Failed to scan student data")
		}
		students = append(students, student)
	}

	// Check for errors after row iteration
	if rows.Err() != nil {
		return nil, common.NewErrorResponse("DB500", "Error while reading students data")
	}

	return students, nil
}

// CreateLink establishes a link between a teacher and a student based on the student's email.
func CreateLink(ctx context.Context, dbPool *pgxpool.Pool, teacherID int, studentEmail string) *common.Response {
	var code, message string

	// Call the database procedure to create a link
	err := dbPool.QueryRow(ctx, "SELECT code, message FROM fn_create_link($1, $2)", teacherID, studentEmail).Scan(&code, &message)
	if err != nil {
		return common.NewErrorResponse("DB500", "Database error occurred while creating link")
	}

	// Check the result from the database procedure and return appropriate response
	if code != "SUCCESS" {
		return common.NewErrorResponse(code, message)
	}

	return common.NewSuccessResponse("Link created successfully", nil)
}
