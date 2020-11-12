package bill

import (
	"fmt"

	creditcard "github.com/durango/go-credit-card"
)

func (bli *billImpl) ValidateCard() error {
	card := creditcard.Card{Number: "4242424242424242", Cvv: "11111", Month: "02", Year: "2016"}
	fmt.Println(card)
	return nil
}
