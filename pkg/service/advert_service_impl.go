package service

import (
	"database/sql"
	"fmt"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/error_message"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/model"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/repository"
	"strings"

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
	if page < 1 {
		return nil, errors.New("page must be >= 1")
	}

	const pageSize = 10
	offset := (page - 1) * pageSize

	var defaultSortField string
	var defaultSortOrder string
	if strings.TrimSpace(sortField) == "" && strings.TrimSpace(sortOrder) == "" {
		defaultSortField = "id"
		defaultSortOrder = "ASC"
	} else {
		switch sortField {
		case "price":
			defaultSortField = "price"
		case "date":
			defaultSortField = "created_at"
		default:
			return nil, errors.New("invalid sort field: must be 'price' or 'date'")
		}
		switch strings.ToLower(sortOrder) {
		case "asc":
			defaultSortOrder = "ASC"
		case "desc":
			defaultSortOrder = "DESC"
		default:
			return nil, errors.New("invalid sort order: must be 'asc' or 'desc'")
		}
	}

	adverts, err := s.advertRepo.List(ctx, pageSize, offset, defaultSortField, defaultSortOrder)
	if err != nil {
		return nil, err
	}

	summaries := make([]AdvertSummary, 0, len(adverts))
	for _, adv := range adverts {
		mainURL, err := s.photoRepo.GetMainPhotoURL(ctx, adv.ID)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, AdvertSummary{
			ID:           adv.ID,
			Name:         adv.Name,
			MainPhotoURL: mainURL,
			Price:        adv.Price,
		})
	}
	return summaries, nil
}

func (s *advertService) Update(ctx context.Context, id int, input UpdateAdvertInput) error {
	advert, err := s.advertRepo.GetByID(ctx, id)
	if err != nil {
		return error_message.ErrAdvertNotFound
	}

	if input.Name != nil {
		advert.Name = *input.Name
	}
	if input.Description != nil {
		advert.Description = *input.Description
	}
	if input.Price != nil {
		if *input.Price <= 0 {
			return error_message.ErrNotPositivePrice
		}
		advert.Price = *input.Price
	}

	if err := s.advertRepo.Update(ctx, advert); err != nil {
		return err
	}

	if input.Photos != nil {
		if err := s.photoRepo.DeleteByAdvertID(ctx, id); err != nil {
			return err
		}
		for idx, url := range *input.Photos {
			photo := model.Photo{
				AdvertID: id,
				URL:      url,
				Position: idx,
			}
			if err := s.photoRepo.Create(ctx, photo); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *advertService) Delete(ctx context.Context, id int) error {
	_, err := s.advertRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_message.ErrAdvertNotFound
		}
		return fmt.Errorf("service.Delete: advertRepo.GetByID (id=%d): %w", id, err)
	}

	if err := s.advertRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return error_message.ErrAdvertNotFound
		}
		return fmt.Errorf("service.Delete: advertRepo.GetByID (id=%d): %w", id, err)
	}
	return nil
}
