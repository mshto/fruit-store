package sales

// GeneralSale GeneralSale
type GeneralSale struct {
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
