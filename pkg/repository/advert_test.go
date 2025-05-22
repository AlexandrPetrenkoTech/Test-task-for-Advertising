package repository

import (
	"Advertising/pkg/model"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPostgresAdvertRepo_CreateAdvert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAdvertRepo(sqlxDB)

	ad := model.Advert{
		Name:        "Test Ad",
		Description: "Test Description",
		Price:       123.45,
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO adverts (name, description, price, created_at)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
	)).
		WithArgs(ad.Name, ad.Description, ad.Price, ad.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	id, err := repo.CreateAdvert(ad)

	assert.NoError(t, err)
	assert.Equal(t, 42, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresAdvertRepo_ListAdverts_AllSortOptions(t *testing.T) {
	// Prepare sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := NewAdvertRepo(sqlxDB)

	// Common test data
	ads := []model.Advert{
		{
			ID:          1,
			Name:        "A",
			Description: "Desc A",
			Price:       100.00,
			CreatedAt:   time.Date(2025, 5, 10, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			Name:        "B",
			Description: "Desc B",
			Price:       200.00,
			CreatedAt:   time.Date(2025, 5, 12, 9, 0, 0, 0, time.UTC),
		},
	}

	// Build expected rows once
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at"})
	for _, ad := range ads {
		rows.AddRow(ad.ID, ad.Name, ad.Description, ad.Price, ad.CreatedAt)
	}

	// Define test cases for each sort combination
	cases := []struct {
		name  string
		field string
		order string
	}{
		{"ByPriceAsc", "price", "ASC"},
		{"ByPriceDesc", "price", "DESC"},
		{"ByDateAsc", "created_at", "ASC"},
		{"ByDateDesc", "created_at", "DESC"},
	}

	limit, offset := 10, 0

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Expect query with dynamic ORDER BY
			query := fmt.Sprintf(
				`SELECT id, name, description, price, created_at
                 FROM adverts
                 ORDER BY %s %s
                 LIMIT \$1 OFFSET \$2`, tc.field, tc.order,
			)
			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(limit, offset).
				WillReturnRows(rows)

			// Execute
			result, err := repo.ListAdverts(limit, offset, tc.field, tc.order)
			assert.NoError(t, err)
			assert.Equal(t, ads, result)
		})
	}

	// Final assertion
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresAdvertRepo_GetAdvertById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewAdvertRepo(sqlxDB)

	expected := model.Advert{
		ID:          42,
		Name:        "Test Ad",
		Description: "This is a test advertisement",
		Price:       199.99,
		CreatedAt:   time.Date(2025, 5, 20, 14, 30, 0, 0, time.UTC),
	}
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at"}).
		AddRow(expected.ID, expected.Name, expected.Description, expected.Price, expected.CreatedAt)

	// Expect the SELECT query
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description, price, created_at
          FROM adverts
         WHERE id = $1`,
	)).
		WithArgs(expected.ID).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetAdvertByID(expected.ID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresAdvertRepo_UpdateAdvert(t *testing.T) {
	// Prepare sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := NewAdvertRepo(sqlxDB)

	// Test data
	updated := model.Advert{
		ID:          42,
		Name:        "Updated Ad",
		Description: "Updated description",
		Price:       250.00,
	}

	// Expect the UPDATE exec
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE adverts
            SET name = $1,
                description = $2,
                price = $3
          WHERE id = $4`,
	)).
		WithArgs(updated.Name, updated.Description, updated.Price, updated.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Execute
	err = repo.UpdateAdvert(updated)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresAdvertRepo_DeleteAdvert(t *testing.T) {
	// Prepare sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := NewAdvertRepo(sqlxDB)

	// Test data
	idToDelete := 42

	// Expect the DELETE exec
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM adverts WHERE id = $1`,
	)).
		WithArgs(idToDelete).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Execute
	err = repo.DeleteAdvert(idToDelete)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
