package repository

import (
	"MentorTools/pkg/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// ConnectDB establishes a connection pool to the auth_db database.
func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	// Load configuration from the config file
	cfg, err := config.LoadConfig("MentorTools/pkg/config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get the specific configuration for auth_db
	dbConfig := cfg.Databases["auth_db"]

	// Create the database URL for auth_db
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.SSLMode,
	)

	// Connect to the database
	pool, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	log.Println("Connected to auth_db successfully.")
	return pool, nil
}

// CloseDB closes the database pool connection.
func CloseDB(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Println("Database connection pool closed.")
	}
}
