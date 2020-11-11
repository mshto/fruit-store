package entity

import (
	"github.com/google/uuid"
)

// UserProduct UserProduct
type UserProduct struct {
	ProductUUID uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      int       `json:"amount"`
	// ID          uuid.UUID `json:"id"`
	// Name        string    `json:"name"`
	// Count       int       `json:"count"`
	// Price       float32   `json:"price"`
}

// GetUserProduct GetUserProduct
type GetUserProduct struct {
	ProductUUID uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float32   `json:"price"`
	Amount      int       `json:"amount"`
	// ID          uuid.UUID `json:"id"`
	// Name        string    `json:"name"`
	// Count       int       `json:"count"`
	// Price       float32   `json:"price"`
}

// UserCart UserCart
type UserCart struct {
	CartProducts []GetUserProduct `json:"products"`
	Total        string           `json:"total"`
}
