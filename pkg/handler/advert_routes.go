package handler

import (
	"Advertising/pkg/service"
	"github.com/labstack/echo/v4"
)

func NewAdvertHandler(e *echo.Echo, svc service.AdvertService) {
	h := &AdvertHandler{advertSvc: svc}

	g := e.Group("/api/adverts")

	g.POST("", h.CreateAdvert)
	g.GET("", h.ListAdverts)
	g.GET("/:id", h.GetAdvertByID)
	g.PUT("/:id", h.UpdateAdvert)
	g.DELETE("/:id", h.DeleteAdvert)
}
