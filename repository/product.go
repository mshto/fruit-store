package repository

import (
	"database/sql"

	"github.com/mshto/fruit-store/backend/entity"
)

// Products Products
type Products interface {
	GetAll() ([]entity.Product, error)
}

// NewProduct NewProduct
func NewProduct(db *sql.DB) Products {
	return &productsImpl{
		db: db,
	}
}

type productsImpl struct {
	db *sql.DB
}

var (
	getAllProducts = `SELECT * FROM products`
)

// GetAll GetAll
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
