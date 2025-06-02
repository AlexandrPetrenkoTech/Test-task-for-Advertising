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

// ====================
// 1. Моки репозиториев с контекстом
// ====================

// MockAdvertRepo реализует мок для вашего repository.AdvertRepo
type MockAdvertRepo struct {
	mock.Mock
}

// Create создаёт новое объявление и возвращает его ID
func (m *MockAdvertRepo) Create(ctx context.Context, ad model.Advert) (int, error) {
	args := m.Called(ctx, ad)
	return args.Int(0), args.Error(1)
}

// List возвращает постраничный список объявлений с сортировкой
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

// GetByID возвращает объявление по ID
func (m *MockAdvertRepo) GetByID(ctx context.Context, id int) (model.Advert, error) {
	args := m.Called(ctx, id)
	// если первый аргумент не nil, приводим его к model.Advert
	if ad, ok := args.Get(0).(model.Advert); ok {
		return ad, args.Error(1)
	}
	return model.Advert{}, args.Error(1)
}

// Update обновляет существующее объявление
func (m *MockAdvertRepo) Update(ctx context.Context, ad model.Advert) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

// Delete удаляет объявление по ID
func (m *MockAdvertRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ====================
// MockPhotoRepo
// ====================

// MockPhotoRepo реализует мок для вашего repository.PhotoRepo
type MockPhotoRepo struct {
	mock.Mock
}

// GetMainPhotoURL возвращает URL «главной» (первой) фотографии объявления
func (m *MockPhotoRepo) GetMainPhotoURL(ctx context.Context, advertID int) (string, error) {
	args := m.Called(ctx, advertID)
	return args.String(0), args.Error(1)
}

// GetAllPhotoURLs возвращает все URL фотографий для данного объявления
func (m *MockPhotoRepo) GetAllPhotoURLs(ctx context.Context, advertID int) ([]string, error) {
	args := m.Called(ctx, advertID)
	if urls, ok := args.Get(0).([]string); ok {
		return urls, args.Error(1)
	}
	return nil, args.Error(1)
}

// Create сохраняет новую запись о фотографии (обычно — в таблицу photos)
func (m *MockPhotoRepo) Create(ctx context.Context, photo model.Photo) error {
	args := m.Called(ctx, photo)
	return args.Error(0)
}

// DeleteByAdvertID удаляет все фотографии, связанные с объявлением advertID
func (m *MockPhotoRepo) DeleteByAdvertID(ctx context.Context, advertID int) error {
	args := m.Called(ctx, advertID)
	return args.Error(0)
}

// ====================
// 2. Вспомогательные данные
// ====================

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

// ====================
// 3. Тесты для AdvertService
// ====================

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

	// Контекст, который будем передавать в мок-методы
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Ожидаем CreateAdvert(ctx, *model.Advert) → вернёт ID = 1
		mockAdRepo.
			On("CreateAdvert", mock.Anything, mock.MatchedBy(func(ad *model.Advert) bool {
				return ad.Name == input.Name && ad.Description == input.Description && ad.Price == input.Price
			})).
			Return(1, nil)

		// Затем ожидаем InsertPhotos(ctx, advertID, urls)
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
		// Предположим, что на ошибке вставки фото сервис возвращает 0
		assert.Equal(t, 0, id)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})
}

func TestAdvertService_GetByID(t *testing.T) {
	mockAdRepo := new(MockAdvertRepo)
	mockPhRepo := new(MockPhotoRepo)
	svc := service.NewAdvertService(mockAdRepo, mockPhRepo)

	const advertID = 42
	adModel := sampleAdvertModel(advertID)
	ctx := context.Background()

	t.Run("SuccessWithoutFields", func(t *testing.T) {
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(adModel, nil)

		mockPhRepo.
			On("GetMainPhotoURL", mock.Anything, advertID).
			Return("http://main_photo", nil)

		out, err := svc.GetByID(ctx, advertID, false)
		assert.NoError(t, err)

		assert.Equal(t, advertID, out.ID)
		assert.Equal(t, adModel.Name, out.Name)
		assert.Equal(t, adModel.Price, out.Price)
		assert.Equal(t, "http://main_photo", out.MainPhotoURL)
		assert.Empty(t, out.Description)
		assert.Empty(t, out.AllPhotosURLs)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})

	t.Run("SuccessWithFields", func(t *testing.T) {
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(adModel, nil)

		mockPhRepo.
			On("GetMainPhotoURL", mock.Anything, advertID).
			Return("http://main_photo", nil)

		mockPhRepo.
			On("GetAllPhotoURLs", mock.Anything, advertID).
			Return([]string{"http://p1", "http://p2"}, nil)

		out, err := svc.GetByID(ctx, advertID, true)
		assert.NoError(t, err)

		assert.Equal(t, advertID, out.ID)
		assert.Equal(t, adModel.Name, out.Name)
		assert.Equal(t, adModel.Price, out.Price)
		assert.Equal(t, "http://main_photo", out.MainPhotoURL)
		assert.Equal(t, adModel.Description, out.Description)
		assert.Equal(t, []string{"http://p1", "http://p2"}, out.AllPhotosURLs)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnGetAdvertByID", func(t *testing.T) {
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(nil, errors.New("not found"))

		_, err := svc.GetByID(ctx, advertID, false)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertNotCalled(t, "GetMainPhotoURL", mock.Anything, mock.Anything)
	})

	t.Run("ErrorOnGetMainPhotoURL", func(t *testing.T) {
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(adModel, nil)

		mockPhRepo.
			On("GetMainPhotoURL", mock.Anything, advertID).
			Return("", errors.New("no photos"))

		_, err := svc.GetByID(ctx, advertID, false)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnGetAllPhotoURLs_WhenFields", func(t *testing.T) {
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(adModel, nil)

		mockPhRepo.
			On("GetMainPhotoURL", mock.Anything, advertID).
			Return("http://main_photo", nil)

		mockPhRepo.
			On("GetAllPhotoURLs", mock.Anything, advertID).
			Return(nil, errors.New("failed to fetch all photos"))

		_, err := svc.GetByID(ctx, advertID, true)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})
}

func TestAdvertService_List(t *testing.T) {
	mockAdRepo := new(MockAdvertRepo)
	mockPhRepo := new(MockPhotoRepo)
	svc := service.NewAdvertService(mockAdRepo, mockPhRepo)

	page := 2
	limit := 10
	offset := (page - 1) * limit
	sortField := "price"
	sortOrder := "asc"
	ctx := context.Background()

	ad1 := &model.Advert{ID: 1, Name: "A", Price: 10.0, CreatedAt: time.Now()}
	ad2 := &model.Advert{ID: 2, Name: "B", Price: 20.0, CreatedAt: time.Now()}
	listFromRepo := []*model.Advert{ad1, ad2}

	t.Run("Success", func(t *testing.T) {
		mockAdRepo.
			On("ListAdverts", mock.Anything, limit, offset, sortField, sortOrder).
			Return(listFromRepo, nil)

		mockPhRepo.
			On("GetMainPhotoURL", mock.Anything, ad1.ID).
			Return("url1", nil)
		mockPhRepo.
			On("GetMainPhotoURL", mock.Anything, ad2.ID).
			Return("url2", nil)

		out, err := svc.List(ctx, page, sortField, sortOrder)
		assert.NoError(t, err)

		expected := []service.AdvertSummary{
			{ID: 1, Name: "A", MainPhotoURL: "url1", Price: 10.0},
			{ID: 2, Name: "B", MainPhotoURL: "url2", Price: 20.0},
		}
		assert.Equal(t, expected, out)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnListAdverts", func(t *testing.T) {
		mockAdRepo.
			On("ListAdverts", mock.Anything, limit, offset, sortField, sortOrder).
			Return(nil, errors.New("db error"))

		_, err := svc.List(ctx, page, sortField, sortOrder)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertNotCalled(t, "GetMainPhotoURL", mock.Anything, mock.Anything)
	})

	t.Run("ErrorOnGetMainPhotoURL", func(t *testing.T) {
		mockAdRepo.
			On("ListAdverts", mock.Anything, limit, offset, sortField, sortOrder).
			Return(listFromRepo, nil)

		mockPhRepo.
			On("GetMainPhotoURL", mock.Anything, ad1.ID).
			Return("", errors.New("no photo"))

		_, err := svc.List(ctx, page, sortField, sortOrder)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})
}

func TestAdvertService_Update(t *testing.T) {
	mockAdRepo := new(MockAdvertRepo)
	mockPhRepo := new(MockPhotoRepo)
	svc := service.NewAdvertService(mockAdRepo, mockPhRepo)

	const advertID = 55
	updateInput := service.UpdateAdvertInput{
		Name:        strPtr("Updated name"),
		Description: strPtr("Updated desc"),
		Photos:      &[]string{"p1", "p2"},
		Price:       floatPtr(11.11),
	}
	ctx := context.Background()

	updatedModel := &model.Advert{
		ID:          advertID,
		Name:        *updateInput.Name,
		Description: *updateInput.Description,
		Price:       *updateInput.Price,
		CreatedAt:   time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		existing := sampleAdvertModel(advertID)
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(existing, nil)

		mockAdRepo.
			On("UpdateAdvert", mock.Anything, updatedModel).
			Return(nil)

		mockPhRepo.
			On("InsertPhotos", mock.Anything, advertID, *updateInput.Photos).
			Return(nil)

		err := svc.Update(ctx, advertID, updateInput)
		assert.NoError(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnGetAdvert", func(t *testing.T) {
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(nil, errors.New("not found"))

		err := svc.Update(ctx, advertID, updateInput)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertNotCalled(t, "InsertPhotos", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ErrorOnUpdateAdvert", func(t *testing.T) {
		existing := sampleAdvertModel(advertID)
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(existing, nil)

		mockAdRepo.
			On("UpdateAdvert", mock.Anything, mock.Anything).
			Return(errors.New("update fail"))

		err := svc.Update(ctx, advertID, updateInput)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertNotCalled(t, "InsertPhotos", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ErrorOnInsertPhotos", func(t *testing.T) {
		existing := sampleAdvertModel(advertID)
		mockAdRepo.
			On("GetAdvertByID", mock.Anything, advertID).
			Return(existing, nil)

		mockAdRepo.
			On("UpdateAdvert", mock.Anything, mock.Anything).
			Return(nil)

		mockPhRepo.
			On("InsertPhotos", mock.Anything, advertID, *updateInput.Photos).
			Return(errors.New("insert photos fail"))

		err := svc.Update(ctx, advertID, updateInput)
		assert.Error(t, err)

		mockAdRepo.AssertExpectations(t)
		mockPhRepo.AssertExpectations(t)
	})
}

func TestAdvertService_Delete(t *testing.T) {
	mockAdRepo := new(MockAdvertRepo)
	mockPhRepo := new(MockPhotoRepo)
	svc := service.NewAdvertService(mockAdRepo, mockPhRepo)

	const advertID = 77
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockPhRepo.
			On("GetAllPhotoURLs", mock.Anything, advertID).
			Return([]string{"a", "b"}, nil)

		mockAdRepo.
			On("DeleteAdvert", mock.Anything, advertID).
			Return(nil)

		err := svc.Delete(ctx, advertID)
		assert.NoError(t, err)

		mockPhRepo.AssertExpectations(t)
		mockAdRepo.AssertExpectations(t)
	})

	t.Run("ErrorOnGetAllPhotoURLs", func(t *testing.T) {
		mockPhRepo.
			On("GetAllPhotoURLs", mock.Anything, advertID).
			Return(nil, errors.New("cannot fetch photos"))

		err := svc.Delete(ctx, advertID)
		assert.Error(t, err)

		mockPhRepo.AssertExpectations(t)
		mockAdRepo.AssertNotCalled(t, "DeleteAdvert", mock.Anything, mock.Anything)
	})

	t.Run("ErrorOnDeleteAdvert", func(t *testing.T) {
		mockPhRepo.
			On("GetAllPhotoURLs", mock.Anything, advertID).
			Return([]string{}, nil)

		mockAdRepo.
			On("DeleteAdvert", mock.Anything, advertID).
			Return(errors.New("delete fail"))

		err := svc.Delete(ctx, advertID)
		assert.Error(t, err)

		mockPhRepo.AssertExpectations(t)
		mockAdRepo.AssertExpectations(t)
	})
}

// ====================
// 4. Утилиты
// ====================

func strPtr(s string) *string     { return &s }
func floatPtr(f float64) *float64 { return &f }
