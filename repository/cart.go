package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/entity"
)

// Cart Cart
type Cart interface {
	GetUserProducts(userUUID uuid.UUID) ([]entity.UserProduct, error)
}

// NewCartProduct NewCartProduct
func NewCartProduct(db *sql.DB) Cart {
	return &productsImpl{
		db: db,
	}
}

type cartImpl struct {
	db *sql.DB
}

var (
	getUserProducts = `SELECT * FROM products WHERE user_id = $1`
)

// GetAll GetAll
func (pri *productsImpl) GetUserProducts(userUUID uuid.UUID) ([]entity.UserProduct, error) {
	products := []entity.UserProduct{}

	rows, err := pri.db.Query(getUserProducts, userUUID)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		p := entity.UserProduct{}
		err := rows.Scan(&p.ID, &p.Name, &p.Count)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	return products, nil
}
