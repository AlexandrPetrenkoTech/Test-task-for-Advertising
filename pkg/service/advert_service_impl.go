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

func (s *advertService) GetByID(id int, fields bool) (AdvertDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (s *advertService) List(page int, sortField, sortOrder string) ([]AdvertSummary, error) {
	//TODO implement me
	panic("implement me")
}

func (s *advertService) Update(id int, input UpdateAdvertInput) error {
	//TODO implement me
	panic("implement me")
}

func (s *advertService) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
