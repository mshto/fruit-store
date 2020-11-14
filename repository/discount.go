package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

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
	getDiscount      = `SELECT id, rule, elements, discount FROM discount WHERE id=$1`
	deleteDiscount   = `DELETE FROM discount WHERE id=$1`
	validateDiscount = "SELECT exists (SELECT id FROM discount WHERE id=$1)"
)

// error
var (
	ErrNotFound = errors.New("not found")
)

func (dsi *discountImpl) GetDiscount(discountID string) (config.GeneralSale, error) {
	sale := config.GeneralSale{}
	var skills []byte

	exists, err := dsi.isRowExist(discountID)
	if err != nil {
		return sale, err
	}
	if !exists {
		return sale, ErrNotFound
	}

	err = dsi.db.QueryRow(getDiscount, discountID).Scan(&sale.ID, &sale.Rule, &skills, &sale.Discount)
	if err != nil {
		return sale, err
	}
	// TODO: add a new struct to unmarshal json
	err = json.Unmarshal(skills, &sale.Elements)
	return sale, err
}

func (dsi *discountImpl) RemoveDiscount(discountID string) error {
	_, err := dsi.db.Exec(deleteDiscount, discountID)
	return err
}

func (dsi *discountImpl) isRowExist(discountID string) (bool, error) {
	var exists bool
	err := dsi.db.QueryRow(validateDiscount, discountID).Scan(&exists)
	return exists, err
}
