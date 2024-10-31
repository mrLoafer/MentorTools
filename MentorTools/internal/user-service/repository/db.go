package repository

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ConnectDB establishes a connection pool to the PostgreSQL database.
func ConnectDB() *pgxpool.Pool {
	// Get the database URL from environment variables or use a default one
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://loafer:Tesla846@localhost:5432/mentor_tools?sslmode=disable"
	}

	// Configure and connect to the pool
	pool, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	log.Println("Connected to the database successfully.")
	return pool
}

// CloseDB closes the database pool connection.
func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
	log.Println("Database connection pool closed.")
}
