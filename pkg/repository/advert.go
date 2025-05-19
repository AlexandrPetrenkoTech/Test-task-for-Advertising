// pkg/repository/advert.go
package repository

import (
	"fmt"

	"Advertising/pkg/model"
	"github.com/jmoiron/sqlx"
)

type AdvertRepo interface {
	// Create a new advert and return its ID
	CreateAdvert(ad model.Advert) (int, error)
	// Retrieve list of adverts with pagination & sorting
	ListAdverts(limit, offset int, sortField, sortOrder string) ([]model.Advert, error)
	// Get single advert by ID
	GetAdvertByID(id int) (model.Advert, error)
	// Update an existing advert
	UpdateAdvert(ad model.Advert) error
	// Delete advert by ID (cascade removes photos)
	DeleteAdvert(id int) error
	// Helpers to fetch photo URLs
	GetMainPhotoURL(advertID int) (string, error)
	GetAllPhotoURLs(advertID int) ([]string, error)
}

type PostgresAdvertRepo struct {
	db *sqlx.DB
}

func NewAdvertRepo(db *sqlx.DB) AdvertRepo {
	return &PostgresAdvertRepo{db: db}
}

func (r *PostgresAdvertRepo) CreateAdvert(ad model.Advert) (int, error) {
	var id int
	err := r.db.QueryRow(
		`INSERT INTO adverts (name, description, price, created_at)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
		ad.Name, ad.Description, ad.Price, ad.CreatedAt,
	).Scan(&id)
	return id, err
}

func (r *PostgresAdvertRepo) ListAdverts(limit, offset int, sortField, sortOrder string) ([]model.Advert, error) {
	var ads []model.Advert
	query := fmt.Sprintf(`
        SELECT id, name, description, price, created_at
          FROM adverts
         ORDER BY %s %s
         LIMIT $1 OFFSET $2`, sortField, sortOrder)
	if err := r.db.Select(&ads, query, limit, offset); err != nil {
		return nil, err
	}
	return ads, nil
}

func (r *PostgresAdvertRepo) GetAdvertByID(id int) (model.Advert, error) {
	var ad model.Advert
	err := r.db.Get(&ad, `
        SELECT id, name, description, price, created_at
          FROM adverts
         WHERE id = $1`, id)
	return ad, err
}

func (r *PostgresAdvertRepo) UpdateAdvert(ad model.Advert) error {
	_, err := r.db.Exec(
		`UPDATE adverts
            SET name = $1,
                description = $2,
                price = $3
          WHERE id = $4`,
		ad.Name, ad.Description, ad.Price, ad.ID,
	)
	return err
}

func (r *PostgresAdvertRepo) DeleteAdvert(id int) error {
	_, err := r.db.Exec(`DELETE FROM adverts WHERE id = $1`, id)
	return err
}

func (r *PostgresAdvertRepo) GetMainPhotoURL(advertID int) (string, error) {
	var url string
	err := r.db.Get(&url, `
        SELECT url
          FROM photos
         WHERE advert_id = $1
      ORDER BY position
         LIMIT 1`, advertID)
	return url, err
}

func (r *PostgresAdvertRepo) GetAllPhotoURLs(advertID int) ([]string, error) {
	var urls []string
	err := r.db.Select(&urls, `
        SELECT url
          FROM photos
         WHERE advert_id = $1
      ORDER BY position`, advertID)
	return urls, err
}
