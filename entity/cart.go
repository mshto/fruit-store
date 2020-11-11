package entity

import (
	"github.com/google/uuid"
)

// UserProduct UserProduct
type UserProduct struct {
	UserID      uuid.UUID `json:"userId"`
	ProductUUID uuid.UUID `json:"productId"`
	Amount      int       `json:"amount"`
	// ID          uuid.UUID `json:"id"`
	// Name        string    `json:"name"`
	// Count       int       `json:"count"`
	// Price       float32   `json:"price"`
}

// GetUserProduct GetUserProduct
type GetUserProduct struct {
	Name   string  `json:"name"`
	Price  float32 `json:"price"`
	Amount int     `json:"amount"`
	// ID          uuid.UUID `json:"id"`
	// Name        string    `json:"name"`
	// Count       int       `json:"count"`
	// Price       float32   `json:"price"`
}

// UserCart UserCart
type UserCart struct {
	Carts []GetUserProduct
	Total string
}
