package bill

import (
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
)

// GeneralSale GeneralSale
type GeneralSale struct {
	Name     string
	Elements map[string]int
	Rules    string
	Discount int
}

// Sales Sales
var Sales = []GeneralSale{
	GeneralSale{
		Elements: map[string]int{
			"Apples": 9,
		},
		Rules:    "more",
		Discount: 10,
	},
	GeneralSale{
		Elements: map[string]int{
			"Pears":   4,
			"Bananas": 2,
		},
		Rules:    "eq",
		Discount: 30,
	},
}

// Bill Bill
type Bill interface {
}

// GeneralSale GeneralSale
type Sale struct {
	Name     string
	Elements map[string]int
	Discount int
}

// New New
func New(cfg *config.Config) Bill {
	return &billImpl{
		cfg: cfg,
	}
}

type billImpl struct {
	cfg *config.Config
}

func (bli *billImpl) GetTotalPrise(products []entity.GetUserProduct) (string, error) {
	// var totalPrise float32

	// element
	bli.validateGeneralSales(&products)
	// return aui.cache.Get(accessUUID)
	return "", nil
}

func (bli *billImpl) validateGeneralSales(products []entity.GetUserProduct) {
	var totalPrise float32
	for _, sale := range Sales {
		var isSale map[string]int
		for productK := range sale.Elements {
			for _, cartProduct := range products {
				if cartProduct.Name != productK {
					continue
				}
			// 	// isSale
			// }
		}
	}
}

// func (ph cartHandler) calculateTotal(products []entity.GetUserProduct) string {
// 	var total float32
// 	for _, prd := range products {
// 		total += float32(prd.Amount) * prd.Price
// 	}

// 	return fmt.Sprintf("%.2f", total)
// }
