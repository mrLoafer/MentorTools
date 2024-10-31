// repository/db.go
package repository

import (
	"MentorTools/pkg/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

var DBPool *pgxpool.Pool

// InitDB initializes a global connection pool to the database.
func InitDB(ctx context.Context) error {
	cfg, err := config.LoadConfig("MentorTools/pkg/config/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	dbConfig := cfg.Databases["auth_db"]
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.SSLMode,
	)

	DBPool, err = pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	log.Println("Connected to auth_db successfully.")
	return nil
}

// CloseDB closes the global database pool connection.
func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection pool closed.")
	}
}
