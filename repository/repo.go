package repository

import (
	"database/sql"
)

// New initializes new repo container for each table entity
func New(db *sql.DB) *Repository {
	return &Repository{
		Product:  NewProduct(db),
		Cart:     NewCartProduct(db),
		Auth:     NewAuth(db),
		Discount: NewDiscount(db),
	}
}

// Repository container for each table entity
type Repository struct {
	Product  Products
	Cart     Cart
	Auth     Auth
	Discount Discount
}
