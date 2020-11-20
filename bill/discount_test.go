package bill

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/mshto/fruit-store/cache"
	redismock "github.com/mshto/fruit-store/cache/mock"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
)

func TestGetDiscountByUser(t *testing.T) {
	type expected struct {
		sale  config.GeneralSale
		isErr bool
	}
	type payload struct {
		cfg       *config.Config
		userUUID  uuid.UUID
		cacheMock func(cacheMock *redismock.MockCache)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Get discount by user with success",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return(`{
						"Rule": "more",
						"Discount": 10
					}`, nil)
				},
			},

			expected: expected{
				sale: config.GeneralSale{
					Rule:     "more",
					Discount: 10,
				},
				isErr: false,
			},
		},
		{
			name: "Get discount by user with failed",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return("", entity.ErrUserNotFound)
				},
			},

			expected: expected{
				sale:  config.GeneralSale{},
				isErr: true,
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

			sale, err := bill.GetDiscountByUser(test.payload.userUUID)
			assert.Equal(t, sale, test.expected.sale)
			if test.expected.isErr {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestSetDiscount(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		cfg       *config.Config
		userUUID  uuid.UUID
		sale      config.GeneralSale
		cacheMock func(cacheMock *redismock.MockCache)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Set discount with success",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				sale:     config.GeneralSale{},
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},

			expected: expected{
				err: nil,
			},
		},
		{
			name: "Set discount with failed",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				sale:     config.GeneralSale{},
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(cache.ErrNotFound)
				},
			},

			expected: expected{
				err: cache.ErrNotFound,
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

			err := bill.SetDiscount(test.payload.userUUID, test.payload.sale)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestRemoveDiscount(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		cfg       *config.Config
		userUUID  uuid.UUID
		cacheMock func(cacheMock *redismock.MockCache)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Remove discount with success",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return(`{
						"Rule": "more",
						"Discount": 10
					}`, nil)
					cacheMock.EXPECT().Del(gomock.Any()).Return(nil)
				},
			},

			expected: expected{
				err: nil,
			},
		},
		{
			name: "Remove discount ErrNotFound with failed",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return("", cache.ErrNotFound)
				},
			},

			expected: expected{
				err: nil,
			},
		},
		{
			name: "Remove discount with failed",
			payload: payload{
				cfg:      &config.Config{},
				userUUID: uuid.New(),
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return("", entity.ErrUserNotFound)
				},
			},

			expected: expected{
				err: entity.ErrUserNotFound,
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

			err := bill.RemoveDiscount(test.payload.userUUID)
			assert.Equal(t, err, test.expected.err)
		})
	}
}
