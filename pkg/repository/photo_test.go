package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPostgresPhotoRepo_GetMainPhotoURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPhotoRepo(sqlxDB)

	expectedURL := "https://example.com/photo1.jpg"

	rows := sqlmock.NewRows([]string{"url"}).
		AddRow(expectedURL)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT url
           FROM photos
          WHERE advert_id = $1
       ORDER BY position
          LIMIT 1`,
	)).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repo.GetMainPhotoURL(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPhotoRepo_GetAllPhotoURLs(t *testing.T) {
	// Prepare sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPhotoRepo(sqlxDB)

	expectedURLs := []string{
		"https://example.com/photo1.jpg",
		"https://example.com/photo2.jpg",
		"https://example.com/photo3.jpg",
	}

	rows := sqlmock.NewRows([]string{"url"})
	for _, url := range expectedURLs {
		rows.AddRow(url)
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT url
           FROM photos
          WHERE advert_id = $1
       ORDER BY position`,
	)).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repo.GetAllPhotoURLs(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedURLs, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
