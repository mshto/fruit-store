package bill

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/config"
)

var (
	pattertn = "%s_discount"
)

func (bli *billImpl) GetDiscountByUser(userUUID uuid.UUID) (config.GeneralSale, error) {
	var sale config.GeneralSale
	fmt.Println(fmt.Sprintf(pattertn, userUUID))
	// return bli.cache.Get(fmt.Sprintf(pattertn, userUUID))
	saleStr, err := bli.cache.Get(fmt.Sprintf(pattertn, userUUID))
	if err != nil {
		return sale, err
	}

	// []byte(serialized)
	err = json.Unmarshal([]byte(saleStr), &sale)
	// if err == nil {
	fmt.Println(sale, err)
	return sale, err
	// }
}

func (bli *billImpl) SetDiscount(userUUID uuid.UUID, sale config.GeneralSale) error {
	fmt.Println(fmt.Sprintf(pattertn, userUUID))
	serialized, err := json.Marshal(sale)
	if err != nil {
		return err
	}

	return bli.cache.Set(fmt.Sprintf(pattertn, userUUID), serialized, time.Duration(bli.cfg.Redis.DiscountTTL)*time.Second)

}

// err = aui.cache.Del(fmt.Sprintf("%s++%s", accessUUID, userUUID))

// at := time.Unix(td.AtExpires, 0)
// rt := time.Unix(td.RtExpires, 0)
// now := time.Now()

// err := aui.cache.Set(td.AccessUUID, userUUID.String(), at.Sub(now))
// if err != nil {
// 	return err
// }
