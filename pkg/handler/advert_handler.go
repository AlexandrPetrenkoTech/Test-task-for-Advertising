package handler

import (
	"Advertising/pkg/error_message"
	"Advertising/pkg/service"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// AdvertHandler отвечает за HTTP-эндпоинты /api/adverts.
type AdvertHandler struct {
	advertSvc service.AdvertService
}

// NewAdvertHandler создаёт новый экземпляр и регистрирует маршруты в Echo.
// Предполагается, что e уже инициализирован (echo.New()).

// CreateAdvert обрабатывает POST /api/adverts
func (h *AdvertHandler) CreateAdvert(c echo.Context) error {
	var req CreateAdvertRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, error_message.ErrBadRequestBody)
	}
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err)
	}

	svcInput := service.CreateAdvertInput{
		Name:        req.Name,
		Description: req.Description,
		Photos:      req.Photos,
		Price:       req.Price,
	}

	newID, err := h.advertSvc.Create(c.Request().Context(), svcInput)
	if err != nil {
		// Обрабатываем все возможные «известные» ошибки
		switch {
		case errors.Is(err, error_message.ErrWrongTitle),
			errors.Is(err, error_message.ErrWrongDescription),
			errors.Is(err, error_message.ErrWrongPhotos),
			errors.Is(err, error_message.ErrNotPositivePrice):
			return SendError(c, http.StatusBadRequest, err)
		default:
			return SendError(c, http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": newID})
}

func (h *AdvertHandler) GetAdvertByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		return SendError(c, http.StatusBadRequest, error_message.ErrWrongAdvertID)
	}

	fields := false
	if fieldParams := c.QueryParam("fields"); fieldParams != "" {
		if f, err := strconv.ParseBool(fieldParams); err != nil {
			return SendError(c, http.StatusBadRequest, error_message.ErrWrongFieldsParam)
		} else {
			fields = f
		}
	}

	adv, err := h.advertSvc.GetByID(c.Request().Context(), id, fields)
	if err != nil {
		if errors.Is(err, error_message.ErrAdvertNotFound) {
			return SendError(c, http.StatusNotFound, err)
		}
		return SendError(c, http.StatusInternalServerError, err)
	}

	response := GetAdvertResponse{
		ID:           adv.ID,
		Name:         adv.Name,
		MainPhotoURL: adv.MainPhotoURL,
		Price:        adv.Price,
	}
	if fields {
		response.Description = &adv.Description
		response.AllPhotosURLs = adv.AllPhotosURLs
	}
	return c.JSON(http.StatusOK, response)
}

// ListAdverts обрабатывает GET /api/adverts?page=&sort=
func (h *AdvertHandler) ListAdverts(c echo.Context) error {
	// 1) Парсим page (по умолчанию 1)
	pageParam := c.QueryParam("page")
	page := 1
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	// 2) Читаем sortParam; если пусто — оставляем sortField, sortOrder == ""
	sortParam := c.QueryParam("sort")
	var sortField, sortOrder string
	if sortParam != "" {
		parts := strings.SplitN(sortParam, "_", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid sort parameter; expected 'price_asc', 'date_desc' и т.д.",
			})
		}
		sortField = parts[0] // "price" или "date" (или невалидная строка)
		sortOrder = parts[1] // "asc" или "desc" (или невалидная строка)
	}
	// Если sortParam == "", то sortField == "" и sortOrder == "" —
	// и сервис сам поставит дефолт “id ASC”.

	// 3) Вызываем сервис, передавая пустые строки, если сортировки нет
	listResp, err := h.advertSvc.List(c.Request().Context(), page, sortField, sortOrder)
	if err != nil {
		// Если например sortField/sortOrder оказались некорректны, сервис вернёт ошибку.
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// 4) Отправляем ответ
	return c.JSON(http.StatusOK, listResp)
}

// UpdateAdvert обрабатывает PUT /api/adverts/:id
func (h *AdvertHandler) UpdateAdvert(c echo.Context) error {

}

// DeleteAdvert обрабатывает DELETE /api/adverts/:id
func (h *AdvertHandler) DeleteAdvert(c echo.Context) error {

}
