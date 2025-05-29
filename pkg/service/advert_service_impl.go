package service

import (
	"Advertising/pkg/model"
	"Advertising/pkg/repository"

	"context"
	"errors"
	"time"
)

type advertService struct {
	advertRepo repository.AdvertRepo
	photoRepo  repository.PhotoRepo
}

func NewAdvertService(ar repository.AdvertRepo, pr repository.PhotoRepo) AdvertService {
	return &advertService{
		advertRepo: ar,
		photoRepo:  pr,
	}
}

func (s *advertService) Create(ctx context.Context, input CreateAdvertInput) (int, error) {
	if input.Name == "" {
		return 0, errors.New("name is required")
	}
	if input.Price <= 0 {
		return 0, errors.New("price must be greater than zero")
	}
	advert := model.Advert{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CreatedAt:   time.Now(),
	}
	advertID, err := s.advertRepo.Create(ctx, advert)
	if err != nil {
		return 0, err
	}
	for idx, url := range input.Photos {
		photo := model.Photo{
			AdvertID: advertID,
			URL:      url,
			Position: idx,
		}
		if err := s.photoRepo.Create(ctx, photo); err != nil {
			return advertID, err
		}
	}
	return advertID, nil
}

func (s *advertService) GetByID(ctx context.Context, id int, fields bool) (AdvertDetail, error) {
	advert, err := s.advertRepo.GetByID(ctx, id)
	if err != nil {
		return AdvertDetail{}, err
	}

	mainURL, err := s.photoRepo.GetMainPhotoURL(ctx, id)
	if err != nil {
		return AdvertDetail{}, err
	}

	summary := AdvertSummary{
		ID:           advert.ID,
		Name:         advert.Name,
		MainPhotoURL: mainURL,
		Price:        advert.Price,
	}

	if !fields {
		return AdvertDetail{AdvertSummary: summary}, nil
	}
	photos, err := s.photoRepo.GetAllPhotoURLs(ctx, id)
	if err != nil {
		return AdvertDetail{}, err
	}
	detail := AdvertDetail{
		AdvertSummary: summary,
		Description:   advert.Description,
		AllPhotosURLs: photos,
	}
	return detail, nil
}

func (s *advertService) List(ctx context.Context, page int, sortField, sortOrder string) ([]AdvertSummary, error) {
	//TODO implement me
	panic("implement me")
}

func (s *advertService) Update(ctx context.Context, id int, input UpdateAdvertInput) error {
	//TODO implement me
	panic("implement me")
}

func (s *advertService) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
