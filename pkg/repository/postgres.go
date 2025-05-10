package repository

import (
	"Advertising/configs"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresDB builds DSN from cfg and returns a connected *sqlx.DB
func NewPostgresDB(cfg *configs.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
	)
	return sqlx.Connect("postgres", dsn)
}
