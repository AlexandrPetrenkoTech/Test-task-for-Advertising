package mocks

import (
	"Advertising/pkg/handler"
	"Advertising/pkg/service"
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockAdvertService implements the AdvertService interface with testify/mock
type MockAdvertService struct {
	mock.Mock
}

func (h *MockAdvertService) Create(
	ctx context.Context,
	input service.CreateAdvertInput,
) (int, error) {
	args := h.Called(ctx, input)
	return args.Int(0), args.Error(1)
}

func (h *MockAdvertService) GetByID(
	ctx context.Context,
	id int,
	fields bool,
) (service.AdvertDetail, error) {
	args := h.Called(ctx, id, fields)
	return args.Get(0).(service.AdvertDetail), args.Error(1)
}

func (h *MockAdvertService) List(
	ctx context.Context,
	page int,
	sortField, sortOrder string,
) ([]service.AdvertSummary, error) {
	args := h.Called(ctx, page, sortField, sortOrder)
	return args.Get(0).([]service.AdvertSummary), args.Error(1)
}

func (h *MockAdvertService) Update(
	ctx context.Context,
	id int,
	input service.UpdateAdvertInput,
) error {
	args := h.Called(ctx, id, input)
	return args.Error(0)
}

func (h *MockAdvertService) Delete(
	ctx context.Context,
	id int,
) error {
	args := h.Called(ctx, id)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	// 1. Set up Echo and mock service
	e := echo.New()
	svc := new(MockAdvertService)
	h := handler.NewAdvertHandler(e, svc)

	// 2. Prepare input data and mock expectations
	input := handler.CreateAdvertRequest{
		Name:        "Sample",
		Description: "Desc",
		Photos:      []string{"http://a"},
		Price:       100,
	}
	svc.On("Create", mock.Anything, input).Return(1, nil).Once()

	// 3. Form the HTTP request with JSON body
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/adverts", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// 4. Invoke the handler method
	err := h.CreateAdvert(ctx)
	assert.NoError(t, err)

	// 5. Check the status and response body
	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp struct {
		ID int `json:"id"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, 1, resp.ID)

	// 6. Ensure all mock expectations were met
	svc.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	e := echo.New()
	svc := new(MockAdvertService)
	h := handler.NewAdvertHandler(e, svc)

	expected := service.AdvertDetail{
		AdvertSummary: service.AdvertSummary{
			ID:           42,
			Name:         "My Ad",
			MainPhotoURL: "http://a",
			Price:        500,
		},
		Description:   "Some desc",
		AllPhotosURLs: []string{"http://a", "http://b"},
	}
	svc.On("GetByID", mock.Anything, 42, true).Return(expected, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/adverts/42", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("42")

	err := h.GetAdvertByID(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	var actual service.AdvertDetail
	_ = json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, expected, actual)
	svc.AssertExpectations(t)
}

func TestList_Success(t *testing.T) {
	e := echo.New()
	svc := new(MockAdvertService)
	h := handler.NewAdvertHandler(e, svc)

	expected := []service.AdvertSummary{
		{
			ID:           1,
			Name:         "First Ad",
			MainPhotoURL: "http://a1",
			Price:        100,
		},
		{
			ID:           2,
			Name:         "Second Ad",
			MainPhotoURL: "http://a2",
			Price:        200,
		},
	}
	svc.On("List", mock.Anything, 1, "price", "asc").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/adverts?page=1&sort=price_asc", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.ListAdverts(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 5. Deserialize and compare
	var actual []service.AdvertSummary
	_ = json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, expected, actual)

	svc.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	e := echo.New()
	svc := new(MockAdvertService)
	h := handler.NewAdvertHandler(e, svc)

	// 2. Prepare the request input data and for the mock
	reqBody := handler.UpdateAdvertRequest{
		Name:        "Updated Name",
		Description: "Updated Desc",
		Photos:      []string{"http://new-photo"},
		Price:       250,
	}
	// Assume that the handler converts UpdateAdvertRequest to service.UpdateAdvertInput
	svcInput := service.UpdateAdvertInput{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Photos:      reqBody.Photos,
		Price:       reqBody.Price,
	}

	// Set up the mock: for any context, id=5 and svcInput, return nil
	svc.
		On("Update", mock.Anything, 5, svcInput).
		Return(nil).
		Once()

	// 3. Form the HTTP PUT request /api/adverts/5
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/adverts/5", bytes.NewReader(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("5")

	// 4. Call the handler
	err := h.UpdateAdvert(ctx)
	assert.NoError(t, err)

	// 5. Check status (204 No Content) and empty body
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())

	// 6. Ensure that the mock service received the expected call
	svc.AssertExpectations(t)
}

func TestDeleteAdvert_Success(t *testing.T) {
	// 1. Set up Echo and mock service
	e := echo.New()
	svc := new(MockAdvertService)
	h := handler.NewAdvertHandler(e, svc)

	// 2. Set up the mock: for any context and id=7 return nil (successful deletion)
	svc.
		On("Delete", mock.Anything, 7).
		Return(nil).
		Once()

	// 3. Form the HTTP DELETE request /api/adverts/7
	req := httptest.NewRequest(http.MethodDelete, "/api/adverts/7", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("7")

	// 4. Call the handler
	err := h.DeleteAdvert(ctx)
	assert.NoError(t, err)

	// 5. Check status (204 No Content) and empty body
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())

	// 6. Ensure that the mock service received the expected call
	svc.AssertExpectations(t)
}
