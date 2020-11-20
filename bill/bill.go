package bill

import (
	"fmt"
	"strconv"

	"github.com/mshto/fruit-store/cache"
	"github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
)

// Bill interface
type Bill interface {
	GetTotalInfo(userUUID uuid.UUID, products []entity.GetUserProduct) (TotalInfo, error)

	GetDiscountByUser(userUUID uuid.UUID) (config.GeneralSale, error)
	SetDiscount(userUUID uuid.UUID, sale config.GeneralSale) error
	RemoveDiscount(userUUID uuid.UUID) error

	ValidateCard(pmt entity.Payment) error
	Pay(userUUID uuid.UUID) error
}

// Sale Sale
type Sale struct {
	Name     string
	Elements map[string]int
	Discount int
}

// TotalInfo TotalInfo
type TotalInfo struct {
	Price   string
	Savings string
	Amount  string
}

// New New
func New(cfg *config.Config, log *logrus.Logger, cache cache.Cache) Bill {
	return &billImpl{
		cfg:   cfg,
		log:   log,
		cache: cache,
	}
}

type billImpl struct {
	cfg   *config.Config
	log   *logrus.Logger
	cache cache.Cache
}

func (bli *billImpl) GetTotalInfo(userUUID uuid.UUID, products []entity.GetUserProduct) (TotalInfo, error) {
	var sales []config.GeneralSale
	userDiscount, err := bli.GetDiscountByUser(userUUID)

	switch {
	case err == cache.ErrNotFound:
	case err != nil:
		bli.log.Errorf("failed to get discount by user, error: %v", err)
	default:
		sales = append(sales, userDiscount)
	}
	sales = append(sales, bli.cfg.Sales...)

	prdMap, priceWithoutSale := bli.getPriceWithoutSale(products)
	salePrds, prd := bli.getProductsWithSale(sales, prdMap)

	totalInfo := bli.getTotalInfo(salePrds, prd, priceWithoutSale)
	return totalInfo, nil
}

func (bli *billImpl) getTotalInfo(salePrds []Result, products map[string]ProductMap, price float32) TotalInfo {
	var totalPrice float32
	var amount int

	for _, salePrd := range salePrds {
		totalPrice = totalPrice + (float32(salePrd.Amount) * salePrd.Price * ((100 - float32(salePrd.Discount)) / 100))
		amount = amount + salePrd.Amount
	}

	for _, product := range products {
		totalPrice = totalPrice + float32(product.Amount)*product.Price
		amount = amount + product.Amount
	}

	return TotalInfo{
		Price:   fmt.Sprintf("%.2f", totalPrice),
		Savings: fmt.Sprintf("%.2f", price-totalPrice),
		Amount:  strconv.Itoa(amount),
	}
}

// Result GeneralSale
type Result struct {
	Name     string
	Price    float32
	Amount   int
	Discount int
}

// ProductMap ProductMap
type ProductMap struct {
	Price  float32
	Amount int
}

func (bli *billImpl) getPriceWithoutSale(products []entity.GetUserProduct) (map[string]ProductMap, float32) {
	var totalPrice float32
	prdMap := map[string]ProductMap{}
	for _, product := range products {
		prdMap[product.Name] = ProductMap{
			Price:  product.Price,
			Amount: product.Amount,
		}
		totalPrice = totalPrice + product.Price*float32(product.Amount)
	}
	return prdMap, totalPrice
}

func (bli *billImpl) getProductsWithSale(sales []config.GeneralSale, products map[string]ProductMap) ([]Result, map[string]ProductMap) {
	results := []Result{}
	for _, sale := range sales {
		var count int
		var isElementsMissed bool
		var isCountUpdated bool

		for productK, productV := range sale.Elements {
			product, ok := products[productK]
			if !ok {
				isElementsMissed = true
				break
			}
			crtCount := product.Amount / productV
			fmt.Println(productK)
			fmt.Println(product.Amount, productV)
			if crtCount < count || count == 0 && !isCountUpdated {
				isCountUpdated = true
				count = crtCount
			}
			fmt.Println(count)
		}

		if count == 0 || isElementsMissed {
			continue
		}

		for productK, productV := range sale.Elements {
			product, ok := products[productK]
			if !ok {
				continue
			}
			result := Result{
				Name:     productK,
				Price:    product.Price,
				Discount: sale.Discount,
			}

			switch sale.Rule {
			case "more":
				result.Amount = product.Amount
				product.Amount = 0
			case "eq":
				result.Amount = count * productV
				product.Amount = product.Amount - result.Amount
			default:
				continue
			}

			products[productK] = product
			results = append(results, result)
		}

	}
	return results, products
}
