package repository

import "github.com/jmoiron/sqlx"

type PhotoRepo interface {
	GetMainPhotoURL(advertID int) (string, error)
	GetAllPhotoURLs(advertID int) ([]string, error)
}

type PostgresPhotoRepo struct {
	db *sqlx.DB
}

func NewPhotoRepo(db *sqlx.DB) PhotoRepo {
	return &PostgresPhotoRepo{db: db}
}

func (r *PostgresPhotoRepo) GetMainPhotoURL(advertID int) (string, error) {
	var url string
	err := r.db.Get(&url, `
        SELECT url
          FROM photos
         WHERE advert_id = $1
      ORDER BY position
         LIMIT 1`, advertID)
	return url, err
}

func (r *PostgresPhotoRepo) GetAllPhotoURLs(advertID int) ([]string, error) {
	var urls []string
	err := r.db.Select(&urls, `
        SELECT url
          FROM photos
         WHERE advert_id = $1
      ORDER BY position`, advertID)
	return urls, err
}
