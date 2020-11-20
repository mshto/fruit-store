package authentication

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/mshto/fruit-store/cache"
	redismock "github.com/mshto/fruit-store/cache/mock"
	"github.com/mshto/fruit-store/config"
)

func TestGetUserUUID(t *testing.T) {
	type expected struct {
		accessUUID string
		err        error
	}
	type payload struct {
		cfg        *config.Config
		accessUUID string
		cacheMock  func(cacheMock *redismock.MockCache)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Get user UUID with success",
			payload: payload{
				cfg:        &config.Config{},
				accessUUID: "test",
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Get(gomock.Any()).Return("test", nil)
				},
			},

			expected: expected{
				accessUUID: "test",
				err:        nil,
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
			auth := New(test.payload.cfg, logger, cache)

			test.payload.cacheMock(cache)

			accessUUID, err := auth.GetUserUUID(test.payload.accessUUID)
			assert.Equal(t, accessUUID, test.expected.accessUUID)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestCreateTokens(t *testing.T) {
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
			name: "Create tokens with success",
			payload: payload{
				cfg: &config.Config{
					Auth: config.Auth{
						AccessSecret:               "accessSecret",
						RefreshSecret:              "refreshSecret",
						AccessSecretAtExpiresInMin: 1,
					},
				},
				userUUID: uuid.New(),
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},

			expected: expected{
				err: nil,
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
			auth := New(test.payload.cfg, logger, cache)

			test.payload.cacheMock(cache)

			token, err := auth.CreateTokens(test.payload.userUUID)
			assert.NotNil(t, token)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type expected struct {
		accessUUID string
		err        error
	}
	type payload struct {
		cfg          *config.Config
		refreshToken string
		cacheMock    func(cacheMock *redismock.MockCache)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Refresh tokens with success",
			payload: payload{
				cfg: &config.Config{
					Auth: config.Auth{
						AccessSecret:               "accessSecret",
						RefreshSecret:              "refreshSecret",
						AccessSecretAtExpiresInMin: 1,
					},
				},
				refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY0OTIzNjUsInJlZnJlc2hfdXVpZCI6IjJmZGVlODY0LWM1MmEtNGM4OS04OGRmLTE0YmVlOTgxZGMyYSsrMWMyYzRjMjQtZTkzMC00OWJiLWFjYzctZGRiNTFlNzM0M2EzIiwidXNlcl9pZCI6IjFjMmM0YzI0LWU5MzAtNDliYi1hY2M3LWRkYjUxZTczNDNhMyJ9.OPC6sT9qJ_jWpNomegpr2OToHyreSGnCbtgpKwdYq_s",
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Del(gomock.Any()).Return(nil)
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},

			expected: expected{
				accessUUID: "test",
				err:        nil,
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

			token, err := bill.RefreshTokens(test.payload.refreshToken)
			assert.NotNil(t, token)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestValidateToken(t *testing.T) {
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
			name: "Validate token with success",
			payload: payload{
				cfg: &config.Config{
					Auth: config.Auth{
						AccessSecret:               "accessSecret",
						RefreshSecret:              "refreshSecret",
						AccessSecretAtExpiresInMin: 1,
					},
				},
				userUUID: uuid.New(),
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
					cacheMock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},

			expected: expected{
				err: nil,
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

			token, err := bill.CreateTokens(test.payload.userUUID)

			newTokens, err := bill.ValidateToken(token.AccessToken)
			assert.NotNil(t, newTokens)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestRemoveTokens(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		cfg         *config.Config
		accessToken string
		userUUID    string
		cacheMock   func(cacheMock *redismock.MockCache)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Validate token with success",
			payload: payload{
				cfg: &config.Config{
					Auth: config.Auth{
						AccessSecret:               "accessSecret",
						RefreshSecret:              "refreshSecret",
						AccessSecretAtExpiresInMin: 1,
					},
				},
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImI5MGNjMGM1LTFhN2ItNDY2Ni05ZTAwLWM0YmYxMTVhNjRkNCIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTYwNTg4ODUwOSwidXNlcl9pZCI6IjE5NjZjNzc1LTI4YzktNDNiOC1hZDc0LTA4ZDI5MDA5ODc3NiJ9.UmjQa_Wj8n-a30oqrfDiV8scty_zQEtYNV9DIesBTOk",
				userUUID:    "userUUID",
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Del(gomock.Any()).Return(nil)
					cacheMock.EXPECT().Del(gomock.Any()).Return(nil)
				},
			},

			expected: expected{
				err: nil,
			},
		},
		{
			name: "Validate token with failed",
			payload: payload{
				cfg: &config.Config{
					Auth: config.Auth{
						AccessSecret:               "accessSecret",
						RefreshSecret:              "refreshSecret",
						AccessSecretAtExpiresInMin: 1,
					},
				},
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImI5MGNjMGM1LTFhN2ItNDY2Ni05ZTAwLWM0YmYxMTVhNjRkNCIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTYwNTg4ODUwOSwidXNlcl9pZCI6IjE5NjZjNzc1LTI4YzktNDNiOC1hZDc0LTA4ZDI5MDA5ODc3NiJ9.UmjQa_Wj8n-a30oqrfDiV8scty_zQEtYNV9DIesBTOk",
				userUUID:    "userUUID",
				cacheMock: func(cacheMock *redismock.MockCache) {
					cacheMock.EXPECT().Del(gomock.Any()).Return(cache.ErrNotFound)
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

			err := bill.RemoveTokens(test.payload.accessToken, test.payload.userUUID)
			assert.Equal(t, err, test.expected.err)
		})
	}
}
