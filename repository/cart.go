package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/entity"
)

// Cart Cart
type Cart interface {
	GetUserProducts(userUUID uuid.UUID) ([]entity.GetUserProduct, error)
	CreateUserProducts(userUUID uuid.UUID, prd entity.UserProduct) error
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
	getUserProducts   = `SELECT users_cart.amount,products.name, products.price FROM users_cart LEFT JOIN products ON users_cart.product_id=products.id AND users_cart.user_id=$1;`
	createUserProduct = `INSERT INTO users_cart (user_id, product_id, amount) VALUES ($1, $2, $3) ON CONFLICT (product_id, user_id) DO UPDATE SET amount=$3 RETURNING user_id`
)

// GetAll GetAll
func (pri *productsImpl) GetUserProducts(userUUID uuid.UUID) ([]entity.GetUserProduct, error) {
	products := []entity.GetUserProduct{}

	rows, err := pri.db.Query(getUserProducts, userUUID)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		p := entity.GetUserProduct{}
		err := rows.Scan(&p.Amount, &p.Name, &p.Price)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetAll GetAll
func (pri *productsImpl) CreateUserProducts(userUUID uuid.UUID, prd entity.UserProduct) error {
	return pri.db.QueryRow(createUserProduct, userUUID, prd.ProductUUID, prd.Amount).Scan(&prd.UserID)
}
