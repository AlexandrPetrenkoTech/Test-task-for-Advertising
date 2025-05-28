package repository

import (
	"Advertising/pkg/model"
	"context"
)

type AdvertRepo interface {
	// Create a new advert and return its ID
	Create(ctx context.Context, ad model.Advert) (int, error)
	// Retrieve list of adverts with pagination & sorting
	List(ctx context.Context, limit, offset int, sortField, sortOrder string) ([]model.Advert, error)
	// Get single advert by ID
	GetByID(ctx context.Context, id int) (model.Advert, error)
	// Update an existing advert
	Update(ctx context.Context, ad model.Advert) error
	// Delete advert by ID (cascade removes photos)
	Delete(ctx context.Context, id int) error
}
