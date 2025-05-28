package repository

type PhotoRepo interface {
	GetMainPhotoURL(advertID int) (string, error)
	GetAllPhotoURLs(advertID int) ([]string, error)
}
