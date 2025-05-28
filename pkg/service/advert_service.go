package service

// CreateAdvertInput содержит данные для создания объявления.
type CreateAdvertInput struct {
	Name        string
	Description string
	Photos      []string
	Price       float64
}

// UpdateAdvertInput содержит поля для частичного обновления объявления.
// Любые из них могут быть nil — тогда соответствующее поле не меняется.
type UpdateAdvertInput struct {
	Name        *string
	Description *string
	Photos      *[]string
	Price       *float64
}

// AdvertSummary — то, что возвращается в списке.
// Поля: ID, название, главная картинка (первый URL) и цена.
type AdvertSummary struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	MainPhotoURL string  `json:"main_photo_url"`
	Price        float64 `json:"price"`
}

// AdvertDetail — полное представление одного объявления.
// Включает AdvertSummary + описание + все ссылки на фото.
type AdvertDetail struct {
	AdvertSummary
	Description   string   `json:"description"`
	AllPhotosURLs []string `json:"all_photos_urls"`
}

// AdvertService описывает бизнес-логику работы с объявлениями.
type AdvertService interface {
	// Create создаёт новое объявление и возвращает его ID.
	Create(input CreateAdvertInput) (int, error)

	// GetByID возвращает объявление по ID.
	// Если fields == true, включает Description и AllPhotosURLs,
	// иначе — только AdvertSummary.
	GetByID(id int, fields bool) (AdvertDetail, error)

	// List возвращает постраничный список объявлений.
	// page — номер страницы (1-based),
	// sortField — "price" или "date",
	// sortOrder — "asc" или "desc".
	List(page int, sortField, sortOrder string) ([]AdvertSummary, error)

	// Update частично обновляет объявление по ID.
	// Использует UpdateAdvertInput для определения изменяемых полей.
	Update(id int, input UpdateAdvertInput) error

	// Delete удаляет объявление по ID.
	Delete(id int) error
}
