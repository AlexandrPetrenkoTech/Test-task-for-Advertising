package model

// Photo represents a single image belonging to an Advert.
// Position defines ordering: 1 = main photo, 2+ = gallery order.
type Photo struct {
	ID       int    `db:"id" json:"id"`
	AdvertID int    `db:"advert_id" json:"advert_id"`
	URL      string `db:"url" json:"url"`
	Position int    `db:"position" json:"position"`
}
