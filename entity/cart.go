package entity

import (
	"github.com/google/uuid"
)

// UserProduct struct
type UserProduct struct {
	ProductUUID uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      int       `json:"amount"`
}

// GetUserProduct struct
type GetUserProduct struct {
	ProductUUID uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float32   `json:"price"`
	Amount      int       `json:"amount"`
}

// UserCart struct
type UserCart struct {
	CartProducts    []GetUserProduct `json:"products"`
	TotalPrice      string           `json:"totalPrice"`
	TotalSavings    string           `json:"totalSavings"`
	Amount          string           `json:"totalAmount"`
	IsDiscountAdded bool             `json:"isDiscountAdded"`
}
