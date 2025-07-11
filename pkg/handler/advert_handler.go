package handler

import (
	"errors"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/error_message"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/service"
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

// CreateAdvert
// @Summary     Create a new advertisement
// @Description Create advertisement with title, description, photos and price
// @Tags        adverts
// @Accept      json
// @Produce     json
// @Param       advert body     handler.CreateAdvertRequest true "Advertisement payload"
// @Success     201    {object} map[string]int           "New advert ID"
// @Failure     400    {object} handler.ErrorResponse
// @Failure     500    {object} handler.ErrorResponse
// @Router      /adverts [post]
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

// GetAdvertByID godoc
// @Summary     Get an advertisement by ID
// @Description Retrieve a single advert detail by its ID
// @Tags        adverts
// @Accept      json
// @Produce     json
// @Param       id    path     int                     true "Advert ID"
// @Success     200   {object} handler.GetAdvertResponse
// @Failure     400   {object} handler.ErrorResponse
// @Failure     404   {object} handler.ErrorResponse
// @Router      /adverts/{id} [get]
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

// ListAdverts godoc
// @Summary     List advertisements
// @Description Get list of adverts with optional pagination and sorting
// @Tags        adverts
// @Accept      json
// @Produce     json
// @Param       page  query    int                     false "Page number"
// @Param       size  query    int                     false "Page size"
// @Param       sort  query    string                  false "Sort by field, e.g. price_asc"
// @Success     200   {array}  handler.GetAdvertResponse
// @Failure     500   {object} handler.ErrorResponse
// @Router      /adverts [get]
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

// UpdateAdvert godoc
// @Summary     Update an advertisement
// @Description Update advertisement fields by ID
// @Tags        adverts
// @Accept      json
// @Produce     json
// @Param       id     path     int                     true "Advert ID"
// @Param       advert body     handler.UpdateAdvertRequest true "Advertisement payload"
// @Success     200    {object} handler.GetAdvertResponse
// @Failure     400    {object} handler.ErrorResponse
// @Failure     404    {object} handler.ErrorResponse
// @Failure     500    {object} handler.ErrorResponse
// @Router      /adverts/{id} [put]
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

// DeleteAdvert godoc
// @Summary     Delete an advertisement
// @Description Delete advertisement identified by its ID
// @Tags        adverts
// @Accept      json
// @Produce     json
// @Param       id path int true "Advert ID"
// @Success     204 {string} string "No content"
// @Failure     400 {object} handler.ErrorResponse
// @Failure     404 {object} handler.ErrorResponse
// @Router      /adverts/{id} [delete]
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
