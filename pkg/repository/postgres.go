package repository

import (
	"Advertising/configs"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgres creates and returns a PostgreSQL connection using the provided configuration
func NewPostgresDB(cfg *configs.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to PostgreSQL: %w", err)
	}

	//Ping the DB to ensure it's available
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping PostgreSQL: %w", err)
	}

	return db, nil
}
