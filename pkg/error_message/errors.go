package error_message

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Error: msg}
}

var (
	ErrWrongPageNumber  = errors.New("wrong page number")
	ErrWrongSortParams  = errors.New("wrong sort params")
	ErrWrongAdvertID    = errors.New("wrong advert id")
	ErrWrongFieldsParam = errors.New("wrong fields param")
	ErrWrongTitle       = errors.New("title must contain from 1 to 200 characters")
	ErrWrongDescription = errors.New("description must contain from 1 to 1000 characters")
	ErrWrongPhotos      = errors.New("advert must contain from 1 to 3 photos")
	ErrNotPositivePrice = errors.New("price must be positive number")
	ErrBadRequestBody   = errors.New("invalid request body")
	ErrAdvertNotFound   = errors.New("advert not found")
)
