package mocks

import (
	"Advertising/pkg/handler"
	"Advertising/pkg/model"
	"bytes"
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

func (m *MockAdvertService) CreateAdvert(input handler.CreateAdvertRequest) (int, error) {
	args := m.Called(input)
	return args.Int(0), args.Error(1)
}

func (m *MockAdvertService) GetAdvertByID(id int, full bool) (model.AdvertDetail, error) {
	args := m.Called(id, full)
	return args.Get(0).(model.AdvertDetail), args.Error(1)
}

func (m *MockAdvertService) ListAdverts(page, pageSize int, sort, order string) ([]model.AdvertSummary, error) {
	args := m.Called(page, pageSize, sort, order)
	return args.Get(0).([]model.AdvertSummary), args.Error(1)
}

func (m *MockAdvertService) UpdateAdvert(id int, input handler.UpdateAdvertRequest) error {
	args := m.Called(id, input)
	return args.Error(0)
}

func (m *MockAdvertService) DeleteAdvert(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateAdvert_Success(t *testing.T) {
	// 1. Настраиваем Echo и мок-сервис
	e := echo.New()
	svc := new(MockAdvertService)
	h := handler.NewAdvertHandler(e, svc)

	// 2. Подготавливаем входные данные и ожидания моков
	input := handler.CreateAdvertRequest{
		Name:        "Sample",
		Description: "Desc",
		Photos:      []string{"http://a"},
		Price:       100,
	}
	svc.On("CreateAdvert", input).Return(1, nil).Once()

	// 3. Формируем HTTP-запрос с JSON-телом
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/adverts", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// 4. Вызываем метод handler
	err := h.CreateAdvert(ctx)
	assert.NoError(t, err)

	// 5. Проверяем статус и тело ответа
	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp struct {
		ID int `json:"id"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, 1, resp.ID)

	// 6. Убеждаемся, что все ожидания моков отработали
	svc.AssertExpectations(t)
}
