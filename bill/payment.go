package bill

import (
	"fmt"
	"strconv"

	creditcard "github.com/durango/go-credit-card"

	"github.com/mshto/fruit-store/entity"
)

const (
	yearPattern = "20%d"
)

func (bli *billImpl) ValidateCard(pmt entity.Payment) error {
	var month, year int
	_, err := fmt.Sscanf(pmt.Expiry, "%d/%d", &month, &year)
	if err != nil {
		return err
	}

	card := creditcard.Card{Number: pmt.CardNumber,
		Cvv:   pmt.Cvc,
		Month: strconv.Itoa(month),
		Year:  fmt.Sprintf(yearPattern, year)}

	return card.Validate()
}
