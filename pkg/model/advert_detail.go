package model

type AdvertDetail struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Price        float64  `json:"price"`
	MainPhotoURL string   `json:"main_photo_url"`
	Description  string   `json:"description"`
	AllPhotoURLs []string `json:"all_photo_ur_ls"`
}
