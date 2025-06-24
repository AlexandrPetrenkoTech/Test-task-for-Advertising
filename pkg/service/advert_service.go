package service

import "context"

// CreateAdvertInput contains data for creating an advert.
type CreateAdvertInput struct {
	Name        string
	Description string
	Photos      []string
	Price       float64
}

// UpdateAdvertInput contains fields for partial advert update.
// Any of them can be nil — in that case, the corresponding field is not changed.
type UpdateAdvertInput struct {
	Name        *string
	Description *string
	Photos      *[]string
	Price       *float64
}

// AdvertSummary represents the data returned in the advert list.
// Fields: ID, name, main photo (first URL), and price.
type AdvertSummary struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	MainPhotoURL string  `json:"main_photo_url"`
	Price        float64 `json:"price"`
}

// AdvertDetail represents a full advert view.
// Includes AdvertSummary + description + all photo URLs.
type AdvertDetail struct {
	AdvertSummary
	Description   string   `json:"description"`
	AllPhotosURLs []string `json:"all_photos_urls"`
}

// AdvertService describes the business logic for working with adverts.
type AdvertService interface {
	// Create creates a new advert and returns its ID.
	Create(ctx context.Context, input CreateAdvertInput) (int, error)

	// GetByID returns an advert by ID.
	// If fields == true, includes Description and AllPhotosURLs,
	// otherwise — only AdvertSummary.
	GetByID(ctx context.Context, id int, fields bool) (AdvertDetail, error)

	// List returns a paginated list of adverts.
	// page — page number (1-based),
	// sortField — "price" or "date",
	// sortOrder — "asc" or "desc".
	List(ctx context.Context, page int, sortField, sortOrder string) ([]AdvertSummary, error)

	// Update partially updates an advert by ID.
	// Uses UpdateAdvertInput to determine which fields to change.
	Update(ctx context.Context, id int, input UpdateAdvertInput) error

	// Delete deletes an advert by ID.
	Delete(ctx context.Context, id int) error
}
