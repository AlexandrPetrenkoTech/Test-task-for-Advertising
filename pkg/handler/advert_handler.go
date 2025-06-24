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

// AdvertHandler is responsible for HTTP endpoints under /api/adverts.
type AdvertHandler struct {
	advertSvc service.AdvertService
}

// NewAdvertHandler creates a new instance and registers routes in Echo.
// It is assumed that e has already been initialized (echo.New()).

// CreateAdvert handles POST /api/adverts
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
		// Handle all possible "known" errors
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

// GetAdvertByID handles GET /api/adverts/:id
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

// ListAdverts handles GET /api/adverts?page=&sort=
func (h *AdvertHandler) ListAdverts(c echo.Context) error {
	// 1) Parse page (default is 1)
	pageParam := c.QueryParam("page")
	page := 1
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	// 2) Read sortParam; if empty — leave sortField and sortOrder as empty strings
	sortParam := c.QueryParam("sort")
	var sortField, sortOrder string
	if sortParam != "" {
		parts := strings.SplitN(sortParam, "_", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid sort parameter; expected 'price_asc', 'date_desc', etc.",
			})
		}
		sortField = parts[0] // "price" or "date" (or invalid string)
		sortOrder = parts[1] // "asc" or "desc" (or invalid string)
	}
	// If sortParam == "", then sortField == "" and sortOrder == "" —
	// and the service will apply the default “id ASC”.

	// 3) Call the service, passing empty strings if no sorting
	listResp, err := h.advertSvc.List(c.Request().Context(), page, sortField, sortOrder)
	if err != nil {
		// For example, if sortField/sortOrder turned out invalid, the service will return an error.
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// 4) Send response
	return c.JSON(http.StatusOK, listResp)
}

// UpdateAdvert handles PUT /api/adverts/:id
func (h *AdvertHandler) UpdateAdvert(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		return SendError(c, http.StatusBadRequest, error_message.ErrWrongAdvertID)
	}
	var req UpdateAdvertRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, error_message.ErrBadRequestBody)
	}

	update := service.UpdateAdvertInput{
		Name:        req.Name,
		Description: req.Description,
		Photos:      req.Photos,
		Price:       req.Price,
	}

	if err := h.advertSvc.Update(c.Request().Context(), id, update); err != nil {
		switch {
		case errors.Is(err, error_message.ErrAdvertNotFound):
			return SendError(c, http.StatusNotFound, err)
		case errors.Is(err, error_message.ErrNotPositivePrice):
			return SendError(c, http.StatusBadRequest, err)
		default:
			return SendError(c, http.StatusInternalServerError, err)
		}
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteAdvert handles DELETE /api/adverts/:id
func (h *AdvertHandler) DeleteAdvert(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		return SendError(c, http.StatusBadRequest, error_message.ErrWrongAdvertID)
	}

	if err := h.advertSvc.Delete(c.Request().Context(), id); err != nil {
		switch {
		case errors.Is(err, error_message.ErrAdvertNotFound):
			return SendError(c, http.StatusNotFound, err)
		default:
			return SendError(c, http.StatusInternalServerError, err)
		}
	}
	return c.NoContent(http.StatusNoContent)
}
