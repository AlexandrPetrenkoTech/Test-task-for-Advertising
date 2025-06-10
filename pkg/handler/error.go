package handler

import (
	"github.com/labstack/echo/v4"
)

// ErrorResponse описывает JSON-ответ с ошибкой.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SendError отправляет клиенту JSON { "error": "<msg>" } с нужным HTTP-статусом.
func SendError(c echo.Context, code int, err error) error {
	return c.JSON(code, ErrorResponse{Error: err.Error()})
}
