package postgres

import (
	"Advertising/pkg/repository"
	"context"
	"fmt"

	"Advertising/pkg/model"
	"github.com/jmoiron/sqlx"
)

type PostgresAdvertRepo struct {
	db *sqlx.DB
}

func NewPostgresAdvertRepo(db *sqlx.DB) repository.AdvertRepo {
	return &PostgresAdvertRepo{db: db}
}

func (r *PostgresAdvertRepo) Create(ctx context.Context, ad model.Advert) (int, error) {
	var id int
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO adverts (name, description, price, created_at)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
		ad.Name, ad.Description, ad.Price, ad.CreatedAt,
	).Scan(&id)
	return id, err
}

func (r *PostgresAdvertRepo) List(ctx context.Context, limit, offset int, sortField, sortOrder string) ([]model.Advert, error) {
	var ads []model.Advert
	query := fmt.Sprintf(`
        SELECT id, name, description, price, created_at
          FROM adverts
         ORDER BY %s %s
         LIMIT $1 OFFSET $2`, sortField, sortOrder)
	if err := r.db.SelectContext(ctx, &ads, query, limit, offset); err != nil {
		return nil, err
	}
	return ads, nil
}

func (r *PostgresAdvertRepo) GetByID(ctx context.Context, id int) (model.Advert, error) {
	var ad model.Advert
	err := r.db.GetContext(ctx, &ad, `
        SELECT id, name, description, price, created_at
          FROM adverts
         WHERE id = $1`, id)
	return ad, err
}

func (r *PostgresAdvertRepo) Update(ctx context.Context, ad model.Advert) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE adverts
            SET name = $1,
                description = $2,
                price = $3
          WHERE id = $4`,
		ad.Name, ad.Description, ad.Price, ad.ID,
	)
	return err
}

func (r *PostgresAdvertRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM adverts WHERE id = $1`, id)
	return err
}
