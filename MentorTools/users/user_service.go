package users

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func UpdateUser(db *pgx.Conn, email, name, role, id string) error {
	query := `UPDATE users SET email=$1, name=$2, role=$3, updated_at=NOW() WHERE id=$4`
	_, err := db.Exec(context.Background(), query, email, name, role, id)
	return err
}

func GetUser(db *pgx.Conn, id string) (string, string, string, error) {
	var email, name, role string
	query := `SELECT email, name, role FROM users WHERE id=$1`
	err := db.QueryRow(context.Background(), query, id).Scan(&email, &name, &role)
	return email, name, role, err
}
