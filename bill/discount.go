package bill

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/config"
)

var (
	pattertn = "%s_discount"
)

// GetDiscountByUser get discount by user
func (bli *billImpl) GetDiscountByUser(userUUID uuid.UUID) (config.GeneralSale, error) {
	var sale config.GeneralSale

	saleStr, err := bli.cache.Get(fmt.Sprintf(pattertn, userUUID))
	if err != nil {
		return sale, err
	}

	err = json.Unmarshal([]byte(saleStr), &sale)
	return sale, err
}

func (bli *billImpl) SetDiscount(userUUID uuid.UUID, sale config.GeneralSale) error {
	serialized, err := json.Marshal(sale)
	if err != nil {
		return err
	}

	return bli.cache.Set(fmt.Sprintf(pattertn, userUUID), serialized, time.Duration(bli.cfg.Redis.DiscountTTL)*time.Second)
}

// RemoveDiscount remove discount
func (bli *billImpl) RemoveDiscount(userUUID uuid.UUID) error {
	_, err := bli.GetDiscountByUser(userUUID)
	switch {
	case err == cache.ErrNotFound:
		return nil
	case err != nil:
		return err

	}
	return bli.cache.Del(fmt.Sprintf(pattertn, userUUID))
}
