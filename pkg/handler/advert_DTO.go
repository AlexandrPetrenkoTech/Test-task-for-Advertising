package handler

// CreateAdvertRequest — payload для POST /api/adverts
type CreateAdvertRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Photos      []string `json:"photos" validate:"required"`
	Price       float64  `json:"price" validate:"required"`
}

// AdvertSummaryResponse — элемент списка GET /api/adverts
type AdvertSummaryResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	MainPhotoURL string  `json:"main_photo_url"`
	Price        float64 `json:"price"`
}

// GetAdvertResponse — ответ GET /api/adverts/:id
type GetAdvertResponse struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	MainPhotoURL  string   `json:"main_photo_url"`
	Price         float64  `json:"price"`
	Description   *string  `json:"description,omitempty"`
	AllPhotosURLs []string `json:"all_photos_urls,omitempty"`
}

// UpdateAdvertRequest — payload для PUT /api/adverts/:id
type UpdateAdvertRequest struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Photos      *[]string `json:"photos,omitempty"`
	Price       *float64  `json:"price,omitempty"`
}
