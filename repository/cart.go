package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/entity"
)

// Cart Cart
type Cart interface {
	GetUserProducts(userUUID uuid.UUID) ([]entity.GetUserProduct, error)
	CreateUserProduct(userUUID uuid.UUID, prd entity.UserProduct) error
	CreateUserProducts(userUUID uuid.UUID, prd entity.UserProduct) error
	RemoveUserProduct(userUUID uuid.UUID, prd entity.UserProduct) error
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

// INSERT INTO users_cart (user_id, product_id, amount) VALUES ('6f047e18-239f-11eb-a734-0242ac150002', '4e558492-239e-11eb-8182-0242ac150002', 1) ON CONFLICT (product_id, user_id) DO UPDATE SET amount=users_cart.amount+1 RETURNING user_id;
var (
	getUserProducts    = `SELECT users_cart.amount, products.id, products.name, products.price FROM users_cart LEFT JOIN products ON users_cart.product_id=products.id AND users_cart.user_id=$1;`
	createUserProducts = `INSERT INTO users_cart (user_id, product_id, amount) VALUES ($1, $2, $3) ON CONFLICT (product_id, user_id) DO UPDATE SET amount=$3 RETURNING user_id`
	createUserProduct  = `INSERT INTO users_cart (user_id, product_id, amount) VALUES ($1, $2, $3) ON CONFLICT (product_id, user_id) DO UPDATE SET amount=users_cart.amount+1 RETURNING user_id`
	deleteUserProduct  = `DELETE FROM users_cart WHERE user_id = $1 AND product_id = $2`
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
		// mshto
		err := rows.Scan(&p.Amount, &p.ProductUUID, &p.Name, &p.Price)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetAll GetAll
func (pri *productsImpl) CreateUserProducts(userUUID uuid.UUID, prd entity.UserProduct) error {
	return pri.db.QueryRow(createUserProducts, userUUID, prd.ProductUUID, prd.Amount).Scan(&prd.UserID)
}

// GetAll GetAll
func (pri *productsImpl) CreateUserProduct(userUUID uuid.UUID, prd entity.UserProduct) error {
	return pri.db.QueryRow(createUserProduct, userUUID, prd.ProductUUID, 1).Scan(&prd.UserID)
}

// GetAll GetAll
func (pri *productsImpl) RemoveUserProduct(userUUID uuid.UUID, prd entity.UserProduct) error {
	_, err := pri.db.Exec(deleteUserProduct, userUUID, prd.ProductUUID)
	return err
}
