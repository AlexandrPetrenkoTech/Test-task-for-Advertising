package repository

import (
	"context"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/model"
)

type PhotoRepo interface {
	GetMainPhotoURL(ctx context.Context, advertID int) (string, error)
	GetAllPhotoURLs(ctx context.Context, advertID int) ([]string, error)
	Create(ctx context.Context, photo model.Photo) error
	DeleteByAdvertID(ctx context.Context, advertID int) error
}
