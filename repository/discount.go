package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/mshto/fruit-store/config"
)

// Discount Discount
type Discount interface {
	GetDiscount(discountID string) (config.GeneralSale, error)
	RemoveDiscount(discountID string) error
}

// NewDiscount NewDiscount
func NewDiscount(db *sql.DB) Discount {
	return &discountImpl{
		db: db,
	}
}

type discountImpl struct {
	db *sql.DB
}

var (
	getDiscount    = `SELECT * FROM discount WHERE id=$1`
	deleteDiscount = `DELETE FROM discount WHERE id=$1`
)

func (dsi *discountImpl) GetDiscount(discountID string) (config.GeneralSale, error) {
	sale := config.GeneralSale{}
	var skills []byte
	err := dsi.db.QueryRow(getDiscount, discountID).Scan(&sale.ID, &sale.Rule, &skills, &sale.Discount)

	// TODO: add a new struct to unmarshal json
	err = json.Unmarshal(skills, &sale.Elements)
	return sale, err
}

func (dsi *discountImpl) RemoveDiscount(discountID string) error {
	_, err := dsi.db.Exec(deleteDiscount, discountID)
	return err
}
