package repository

import "database/sql"

// New initializes new repo container for each table entity
func New(db *sql.DB) *Repository {
	return &Repository{
		Product: NewProduct(db),
		Cart:    NewCartProduct(db),
		Auth:    NewAuth(db),
	}
}

// Repository Repository
type Repository struct {
	Product Products
	Cart    Cart
	Auth    Auth
}
