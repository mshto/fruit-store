package entity

// Payment Payment
type Payment struct {
	CardNumber string `json:"number"`
	Expiry     string `json:"expiry"`
	Name       string `json:"name"`
	Cvc        string `json:"cvc"`
}
