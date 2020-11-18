package repository

import (
	"database/sql"

	"github.com/mshto/fruit-store/entity"
)

// Products interface
type Products interface {
	GetAll() ([]entity.Product, error)
}

// NewProduct generate a new product
func NewProduct(db *sql.DB) Products {
	return &productsImpl{
		db: db,
	}
}

type productsImpl struct {
	db *sql.DB
}

var (
	getAllProducts = `SELECT id, name, price, created_at FROM products`
)

// GetAll products
func (pri *productsImpl) GetAll() ([]entity.Product, error) {
	products := []entity.Product{}

	rows, err := pri.db.Query(getAllProducts)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		p := entity.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	return products, nil
}
