package repository

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/entity"
)

// Cart interface
type Cart interface {
	GetUserProducts(userUUID uuid.UUID) ([]entity.GetUserProduct, error)
	CreateUserProduct(userUUID, productUUID uuid.UUID) error
	CreateUserProducts(userUUID uuid.UUID, prd entity.UserProduct) error
	RemoveUserProducts(userUUID uuid.UUID) error
	RemoveUserProduct(userUUID, productUUID uuid.UUID) error
}

// NewCartProduct generate a new cart product
func NewCartProduct(db *sql.DB) Cart {
	return &productsImpl{
		db: db,
	}
}

type cartImpl struct {
	db *sql.DB
}

var (
	getUserProducts    = `SELECT users_cart.amount, products.id, products.name, products.price FROM users_cart LEFT JOIN products ON users_cart.product_id=products.id AND users_cart.user_id=$1;`
	createUserProducts = `INSERT INTO users_cart (user_id, product_id, amount) VALUES ($1, $2, $3) ON CONFLICT (product_id, user_id) DO UPDATE SET amount=$3 RETURNING user_id`
	createUserProduct  = `INSERT INTO users_cart (user_id, product_id, amount) VALUES ($1, $2, $3) ON CONFLICT (product_id, user_id) DO UPDATE SET amount=users_cart.amount+1 RETURNING user_id`
	deleteUserProducts = `DELETE FROM users_cart WHERE user_id = $1`
	deleteUserProduct  = `DELETE FROM users_cart WHERE user_id = $1 AND product_id = $2`
)

// GetUserProducts get user products
func (pri *productsImpl) GetUserProducts(userUUID uuid.UUID) ([]entity.GetUserProduct, error) {
	products := []entity.GetUserProduct{}

	rows, err := pri.db.Query(getUserProducts, userUUID)
	if err != nil {
		return products, err
	}

	defer rows.Close()

	for rows.Next() {
		p := entity.GetUserProduct{}
		err := rows.Scan(&p.Amount, &p.ProductUUID, &p.Name, &p.Price)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}
	return products, nil
}

// CreateUserProducts create user products
func (pri *productsImpl) CreateUserProducts(userUUID uuid.UUID, prd entity.UserProduct) error {
	return pri.db.QueryRow(createUserProducts, userUUID, prd.ProductUUID, prd.Amount).Scan(&prd.UserID)
}

// CreateUserProduct create user products
func (pri *productsImpl) CreateUserProduct(userUUID, productUUID uuid.UUID) error {
	return pri.db.QueryRow(createUserProduct, userUUID, productUUID, 1).Scan(&userUUID)
}

// RemoveUserProducts remove user products
func (pri *productsImpl) RemoveUserProducts(userUUID uuid.UUID) error {
	_, err := pri.db.Exec(deleteUserProducts, userUUID)
	return err
}

// RemoveUserProduct remove user product
func (pri *productsImpl) RemoveUserProduct(userUUID, productUUID uuid.UUID) error {
	_, err := pri.db.Exec(deleteUserProduct, userUUID, productUUID)
	return err
}
