package handler

import (
	"github.com/labstack/echo/v4"
)

// ErrorResponse describes a JSON error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SendError sends to the client a JSON { "error": "<msg>" } with the specified HTTP status code.
func SendError(c echo.Context, code int, err error) error {
	return c.JSON(code, ErrorResponse{Error: err.Error()})
}
