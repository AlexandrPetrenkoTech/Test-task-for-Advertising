package postgres

import (
	"Advertising/pkg/model"
	"Advertising/pkg/repository"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresPhotoRepo struct {
	db *sqlx.DB
}

func NewPostgresPhotoRepo(db *sqlx.DB) repository.PhotoRepo {
	return &PostgresPhotoRepo{db: db}
}

func (r *PostgresPhotoRepo) Create(ctx context.Context, photo model.Photo) error {
	query := `
        INSERT INTO photos (advert_id, url, position)
        VALUES ($1, $2, $3)
    `
	_, err := r.db.ExecContext(ctx, query, photo.AdvertID, photo.URL, photo.Position)
	return err
}

func (r *PostgresPhotoRepo) GetMainPhotoURL(ctx context.Context, advertID int) (string, error) {
	var url string
	err := r.db.GetContext(
		ctx, &url,
		`
        SELECT url
          FROM photos
         WHERE advert_id = $1
      ORDER BY position
         LIMIT 1`, advertID,
	)
	return url, err
}

func (r *PostgresPhotoRepo) GetAllPhotoURLs(ctx context.Context, advertID int) ([]string, error) {
	var urls []string
	err := r.db.SelectContext(
		ctx, &urls,
		`
        SELECT url
          FROM photos
         WHERE advert_id = $1
      ORDER BY position`, advertID,
	)
	return urls, err
}

func (r *PostgresPhotoRepo) DeleteByAdvertID(ctx context.Context, advertID int) error {
	query := `
        DELETE
          FROM photos
         WHERE advert_id = $1
    `
	// ExecContext returns a Result and an error; we only care about the error
	if _, err := r.db.ExecContext(ctx, query, advertID); err != nil {
		return fmt.Errorf("failed to delete photos for advert %d: %w", advertID, err)
	}
	return nil
}
