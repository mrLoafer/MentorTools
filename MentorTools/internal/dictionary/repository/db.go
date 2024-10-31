package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

func ConnectDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://loafer:Tesla846@localhost:5432/mentor_tools?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return conn
}
