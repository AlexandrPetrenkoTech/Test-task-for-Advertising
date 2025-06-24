package service_test

import (
	"Advertising/pkg/model"
	"Advertising/pkg/service"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdvertRepo implements a mock for repository.AdvertRepo
type MockAdvertRepo struct {
	mock.Mock
}

// Create creates a new advert and returns its ID
func (m *MockAdvertRepo) Create(ctx context.Context, ad model.Advert) (int, error) {
	args := m.Called(ctx, ad)
	return args.Int(0), args.Error(1)
}

// List returns a paginated list of adverts with sorting
func (m *MockAdvertRepo) List(
	ctx context.Context,
	limit, offset int,
	sortField, sortOrder string,
) ([]model.Advert, error) {
	args := m.Called(ctx, limit, offset, sortField, sortOrder)
	if stored := args.Get(0); stored != nil {
		return stored.([]model.Advert), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetByID returns an advert by ID
func (m *MockAdvertRepo) GetByID(ctx context.Context, id int) (model.Advert, error) {
	args := m.Called(ctx, id)
	// if the first argument is not nil, cast it to model.Advert
	if ad, ok := args.Get(0).(model.Advert); ok {
		return ad, args.Error(1)
	}
	return model.Advert{}, args.Error(1)
}

// Update updates an existing advert
func (m *MockAdvertRepo) Update(ctx context.Context, ad model.Advert) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

// Delete deletes an advert by ID
func (m *MockAdvertRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockPhotoRepo implements a mock for repository.PhotoRepo
type MockPhotoRepo struct {
	mock.Mock
}

// GetMainPhotoURL returns the URL of the "main" (first) photo of the advert
func (m *MockPhotoRepo) GetMainPhotoURL(ctx context.Context, advertID int) (string, error) {
	args := m.Called(ctx, advertID)
	return args.String(0), args.Error(1)
}

// GetAllPhotoURLs returns all photo URLs for the given advert
func (m *MockPhotoRepo) GetAllPhotoURLs(ctx context.Context, advertID int) ([]string, error) {
	args := m.Called(ctx, advertID)
	if urls, ok := args.Get(0).([]string); ok {
		return urls, args.Error(1)
	}
	return nil, args.Error(1)
}

// Create saves a new photo record (usually into the photos table)
func (m *MockPhotoRepo) Create(ctx context.Context, photo model.Photo) error {
	args := m.Called(ctx, photo)
	return args.Error(0)
}

// DeleteByAdvertID deletes all photos associated with the given advert ID
func (m *MockPhotoRepo) DeleteByAdvertID(ctx context.Context, advertID int) error {
	args := m.Called(ctx, advertID)
	return args.Error(0)
}

func sampleAdvertModel(id int) *model.Advert {
	return &model.Advert{
		ID:          id,
		Name:        "Test name",
		Description: "Test description",
		Price:       123.45,
		CreatedAt:   time.Now(),
	}
}

func samplePhotos() []string {
	return []string{"http://img1", "http://img2", "http://img3"}
}

func TestAdvertService_Create(t *testing.T) {
	mockAdRepo := new(MockAdvertRepo)
	mockPhRepo := new(MockPhotoRepo)
	svc := service.NewAdvertService(mockAdRepo, mockPhRepo)

	input := service.CreateAdvertInput{
		Name:        "New Ad",
		Description: "Desc",
		Photos:      samplePhotos(),
		Price:       99.99,
	}

	// Context passed to mock methods
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Expect CreateAdvert(ctx, *model.Advert) → returns ID = 1
		mockAdRepo.
			On("CreateAdvert", mock.Anything, mock.MatchedBy(func(ad *model.Advert) bool {
				return ad.Name == input.Name && ad.Description == input.Description && ad.Price == input.Price
			})).
			Return(1, nil)

		// Then expect InsertPhotos(ctx, advertID, urls)
		mockPhRepo.
			On("InsertPhotos", mock.Anything, 1, input.Photos).
			Return(nil)

		id, err := svc.Create(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnCreateAdvert", func(t *testing.T) {
		mockAdRepo.
			On("CreateAdvert", mock.Anything, mock.Anything).
			Return(0, errors.New("db error"))

		id, err := svc.Create(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, 0, id)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertNotCalled(t, "InsertPhotos", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ErrorOnInsertPhotos", func(t *testing.T) {
		mockAdRepo.
			On("CreateAdvert", mock.Anything, mock.Anything).
			Return(2, nil)

		mockPhRepo.
			On("InsertPhotos", mock.Anything, 2, input.Photos).
			Return(errors.New("photo insert error"))

		id, err := svc.Create(ctx, input)
		assert.Error(t, err)
		// Assume that on photo insertion error, service returns 0
		assert.Equal(t, 0, id)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})
}

func strPtr(s string) *string     { return &s }
func floatPtr(f float64) *float64 { return &f }
