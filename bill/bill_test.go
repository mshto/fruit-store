package bill

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	redismock "github.com/mshto/fruit-store/cache/mock"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
)

func TestGetTotalInfo(t *testing.T) {
	type expected struct {
		total TotalInfo
		isErr bool
	}
	type payload struct {
		cfg       *config.Config
		userUUID  uuid.UUID
		products  []entity.GetUserProduct
		cacheMock func(cacheMock *redismock.MockCache)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Get total info user sale with success",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				products: []entity.GetUserProduct{
					{
						Name:   "Apples",
						Price:  100,
						Amount: 1,
					},
				},
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return(`{
						"Elements": {
							"Apples": 1
						},
						"Rule": "more",
						"Discount": 10
					}`, nil)
				},
			},

			expected: expected{
				total: TotalInfo{
					Price:   "90.00",
					Savings: "10.00",
					Amount:  "1",
				},
				isErr: false,
			},
		},
		{
			name: "Get total info general sale with success",
			payload: payload{
				cfg: &config.Config{
					Sales: []config.GeneralSale{
						{
							Elements: map[string]int{
								"Apples": 1,
							},
							Rule:     "eq",
							Discount: 10,
						},
						{
							Elements: map[string]int{
								"New": 1,
							},
							Rule:     "new",
							Discount: 10,
						},
					},
				},
				userUUID: uuid.New(),
				products: []entity.GetUserProduct{
					{
						Name:   "Apples",
						Price:  100,
						Amount: 1,
					},
				},
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return("", nil)
				},
			},

			expected: expected{
				total: TotalInfo{
					Price:   "90.00",
					Savings: "10.00",
					Amount:  "1",
				},
				isErr: false,
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

			test.payload.cacheMock(cache)

			total, err := bill.GetTotalInfo(test.payload.userUUID, test.payload.products)
			assert.Equal(t, total, test.expected.total)
			if test.expected.isErr {
				assert.NotNil(t, err)
			}
		})
	}
}
