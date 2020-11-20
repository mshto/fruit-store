package bill

import (
	"testing"

	"github.com/golang/mock/gomock"
	redismock "github.com/mshto/fruit-store/cache/mock"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
)

func TestValidateCard(t *testing.T) {
	type expected struct {
		validateErr func(t *testing.T, err error)
	}
	type payload struct {
		cfg     *config.Config
		payment entity.Payment
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Validate card with success",
			payload: payload{
				cfg: &config.Config{},
				payment: entity.Payment{
					CardNumber: "4916527199683696",
					Expiry:     "11/99",
					Name:       "test",
					Cvc:        "123",
				},
			},

			expected: expected{
				validateErr: func(t *testing.T, err error) {
					assert.Nil(t, err)
				},
			},
		},
		{
			name: "Validate card with failed",
			payload: payload{
				cfg: &config.Config{},
				payment: entity.Payment{
					CardNumber: "4916527199683696",
					Expiry:     "invalid",
					Name:       "test",
					Cvc:        "123",
				},
			},

			expected: expected{
				validateErr: func(t *testing.T, err error) {
					assert.NotNil(t, err)
				},
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			logger, _ := loggermock.NewNullLogger()
			cache := redismock.NewMockCache(mockCtrl)
			bill := New(test.payload.cfg, logger, cache)

			err := bill.ValidateCard(test.payload.payment)
			test.expected.validateErr(t, err)
		})
	}
}
