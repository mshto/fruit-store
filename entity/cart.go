package entity

import (
	"github.com/google/uuid"
)

// UserProduct UserProduct
type UserProduct struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Price float32   `json:"price"`
}
