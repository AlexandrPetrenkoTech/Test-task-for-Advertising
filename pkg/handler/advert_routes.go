package handler

import (
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/service"
	"github.com/labstack/echo/v4"
)

// NewAdvertHandler registers advert routes with Swagger annotations
func NewAdvertHandler(e *echo.Echo, svc service.AdvertService) *AdvertHandler {
	h := &AdvertHandler{advertSvc: svc}

	// Advert group
	g := e.Group("/api/adverts")

	g.POST("", h.CreateAdvert)
	g.GET("", h.ListAdverts)
	g.GET("/:id", h.GetAdvertByID)
	g.PUT("/:id", h.UpdateAdvert)
	g.DELETE("/:id", h.DeleteAdvert)

	return h
}
