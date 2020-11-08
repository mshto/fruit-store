package entity

import (
	"time"

	"github.com/google/uuid"
)

// Product Product
type Product struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}
