package bill

import (
	"fmt"
	"strconv"

	"github.com/mshto/fruit-store/cache"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
)

// Bill Bill
type Bill interface {
	GetTotalInfo(userUUID uuid.UUID, products []entity.GetUserProduct) (TotalInfo, error)

	GetDiscountByUser(userUUID uuid.UUID) (config.GeneralSale, error)
	SetDiscount(userUUID uuid.UUID, sale config.GeneralSale) error
	// GetDiscount(userUUID) (TotalInfo, error)
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
func New(cfg *config.Config, cache *cache.Cache) Bill {
	return &billImpl{
		cfg:   cfg,
		cache: cache,
	}
}

type billImpl struct {
	cfg   *config.Config
	cache *cache.Cache
}

func (bli *billImpl) GetTotalInfo(userUUID uuid.UUID, products []entity.GetUserProduct) (TotalInfo, error) {
	// var totalPrise float32
	var sales []config.GeneralSale
	// get user sales
	// GetDiscountByUser
	userDiscount, err := bli.GetDiscountByUser(userUUID)
	switch {
	case err == cache.ErrNotFound:
	case err != nil:
		// log error here
	default:
		sales = append(sales, userDiscount)
	}
	sales = append(sales, bli.cfg.Sales...)
	fmt.Println("mshto final", sales)
	prdMap, priceWithoutSale := bli.getPriceWithoutSale(products)
	salePrds, prd := bli.getProductsWithSale(sales, prdMap)

	totalInfo := bli.getTotalInfo(salePrds, prd, priceWithoutSale)
	// totalInfo.PriceWithoutSale = priceWithoutSale
	// fmt.Sprintf("%.2f", totalPrice)
	return totalInfo, nil
}

func (bli *billImpl) getTotalInfo(salePrds []Result, products map[string]ProductMap, price float32) TotalInfo {
	var totalPrice float32
	var amount int
	fmt.Println("result sale", salePrds)
	for _, salePrd := range salePrds {
		totalPrice = totalPrice + (float32(salePrd.Amount) * salePrd.Price * ((100 - float32(salePrd.Discount)) / 100))
		amount = amount + salePrd.Amount
	}
	fmt.Println("result price", totalPrice)
	fmt.Println("result total", products)
	for _, product := range products {
		totalPrice = totalPrice + float32(product.Amount)*product.Price
		amount = amount + product.Amount
	}
	fmt.Println(totalPrice)
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

		for productK, productV := range sale.Elements {
			product, ok := products[productK]
			if !ok {
				isElementsMissed = true
				break
			}
			crtCount := product.Amount / productV
			fmt.Println(productK)
			fmt.Println(product.Amount, productV)
			if crtCount < count || count == 0 {
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

// Golang playground
// package main

// import (
// 	"fmt"
// )

// func main() {
// 	products := getProducts()
// 	prdMap:=convertListToMap(products)
// 	fmt.Println("products",products)
// 	res, prd:= validateGeneralSales(prdMap)
// 	fmt.Println("results", res)
// 	fmt.Println("new products", prd)
// }

// // Result GeneralSale
// type Result struct {
// 	Name string
// 	Price       float32
// 	Amount      int
// 	Discount int
// }

// // ProductMap ProductMap
// type ProductMap struct {
// 	Price  float32
// 	Amount int
// }

// func convertListToMap(products []GetUserProduct) map[string]ProductMap {
// 	prdMap := map[string]ProductMap{}
// 	for _, product := range products {
// 		prdMap[product.Name] = ProductMap{
// 			Price:  product.Price,
// 			Amount: product.Amount,
// 		}
// 	}
// 	return prdMap
// }

// func validateGeneralSales(products map[string]ProductMap) ([]Result, map[string]ProductMap) {
// 	// var totalPrise float32
// 	results := []Result{}
// 	for _, sale := range Sales {
// 		// var isSale int
// 		var count int

// 		for productK, productV := range sale.Elements {
// 			product, ok := products[productK]
// 			if !ok {
// 				continue
// 			}
// 			crtCount := product.Amount / productV

// 			if crtCount < count || count == 0 {
// 				count = crtCount
// 			}
// 			//for _, cartProduct := range products {
// 			//	if cartProduct.Name != productK {
// 			//		continue
// 			//	}
// 			//	crtCount := cartProduct.Amount / productV

// 			//	if crtCount < count || count == 0 {
// 			//		count = crtCount
// 			//	}
// 			//}
// 		}

// 		if count == 0 {
// 			continue
// 		}

// 		for productK, productV := range sale.Elements {
// 			//for prdKey, cartProduct := range products {
// 				//if cartProduct.Name != productK {
// 			//		continue
// 			//	}
// 				product, ok := products[productK]
// 				if !ok {
// 					continue
// 				}
// 				result:= Result{
// 					Name: productK,
// 					Price: product.Price,
// 					Discount: sale.Discount,
// 				}

// 				switch sale.Rules {
// 				case "more":
// 					result.Amount = product.Amount
// 					product.Amount = 0
// 				case "eq":
// 					result.Amount = count * productV
// 					product.Amount = product.Amount - result.Amount
// 				default:
// 				fmt.Println("here 1")
// 					continue
// 				}
// 				fmt.Println("here 1")
// 				//if sale.Rules == "more" {
// 				 //  result.Amount = cartProduct.Amount
// 				//	fmt.Println("result.Amount", result.Amount)
// 				   //products[prdKey].Amount = 0
// 				  // 	fmt.Println("cartProduct.Amount", cartProduct.Amount)
// 				//}
// 				//if sale.Rules == "eq" {
// 				   //result.Amount = count * productV
// 				  // products[prdKey].Amount = cartProduct.Amount - result.Amount
// 				//}
// 				products[productK] = product
// 				results = append(results,result)

// 		//	}
// 		}

// 	}
// 	return results, products
// }

// func getProducts() []GetUserProduct{
// 	return []GetUserProduct{
// 		GetUserProduct{
// 			Name: "Apples",
// 			Price: 20,
// 			Amount: 10,
// 		},
// 		GetUserProduct{
// 			Name: "Pears",
// 			Price: 20,
// 			Amount: 10,
// 		},
// 		GetUserProduct{
// 			Name: "Bananas",
// 			Price: 20,
// 			Amount: 10,
// 		},
// 		GetUserProduct{
// 			Name: "Bananase",
// 			Price: 20,
// 			Amount: 10,
// 		},

// 	}
// }

// // GetUserProduct GetUserProduct
// type GetUserProduct struct {
// 	//ProductUUID uuid.UUID `json:"id"`
// 	Name        string    `json:"name"`
// 	Price       float32   `json:"price"`
// 	Amount      int       `json:"amount"`
// 	// ID          uuid.UUID `json:"id"`
// 	// Name        string    `json:"name"`
// 	// Count       int       `json:"count"`
// 	// Price       float32   `json:"price"`
// }

// // GeneralSale GeneralSale
// type GeneralSale struct {
// 	Name     string
// 	Elements map[string]int
// 	Rules    string
// 	Discount int
// }

// // Sales Sales
// var Sales = []GeneralSale{
// 	GeneralSale{
// 		Elements: map[string]int{
// 			"Apples": 9,
// 		},
// 		Rules:    "more",
// 		Discount: 30,
// 	},
// 	GeneralSale{
// 		Elements: map[string]int{
// 			"Pears":   4,
// 			"Bananas": 2,
// 		},
// 		Rules:    "eq",
// 		Discount: 30,
// 	},
// 		GeneralSale{
// 		Elements: map[string]int{
// 			"Fails":   4,
// 		},
// 		Rules:    "eq",
// 		Discount: 30,
// 	},
// }
