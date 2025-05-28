package repository

import (
	"Advertising/pkg/model"
	"context"
)

type AdvertRepo interface {
	// Create a new advert and return its ID
	CreateAdvert(ctx context.Context, ad model.Advert) (int, error)
	// Retrieve list of adverts with pagination & sorting
	ListAdverts(ctx context.Context, limit, offset int, sortField, sortOrder string) ([]model.Advert, error)
	// Get single advert by ID
	GetAdvertByID(ctx context.Context, id int) (model.Advert, error)
	// Update an existing advert
	UpdateAdvert(ctx context.Context, ad model.Advert) error
	// Delete advert by ID (cascade removes photos)
	DeleteAdvert(ctx context.Context, id int) error
}
