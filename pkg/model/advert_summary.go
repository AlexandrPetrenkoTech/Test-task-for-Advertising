package model

type AdvertSummary struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	MainPhotoURL string  `json:"main_photo_url"`
}
