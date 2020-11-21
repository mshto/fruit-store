package auth

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	authmock "github.com/mshto/fruit-store/authentication/mock"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
	repomock "github.com/mshto/fruit-store/repository/mock"
	"github.com/mshto/fruit-store/web/middleware"
)

func TestSignup(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		body     []byte
		repoMock func(repoMock *repomock.MockAuth)
		authMock func(authMock *authmock.MockAuth)
	}
	type expected struct {
		code int
		body string
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Sign up with success",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"username":"test","password":"password","passwordRepeat":"password"}`),
				repoMock: func(repoMock *repomock.MockAuth) {
					repoMock.EXPECT().Signup(gomock.Any()).Return(nil)
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
			},
			expected: expected{
				code: http.StatusCreated,
				body: `{}`,
			},
		},
		{
			name: "Sign up invalid body with fail",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`invalid`),
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid character 'i' looking for beginning of value"}`,
			},
		},
		{
			name: "Sign up invalid password with fail",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"username":"test","password":"password","passwordRepeat":""}`),
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
			},
			expected: expected{
				code: http.StatusNotFound,
				body: `{"error":"passwords aren't equal"}`,
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

			authRepo := repomock.NewMockAuth(mockCtrl)
			test.payload.repoMock(authRepo)

			auth := authmock.NewMockAuth(mockCtrl)
			test.payload.authMock(auth)

			req, err := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer(test.payload.body))
			if err != nil {
				t.Error("failed to create request")
			}
			rw := httptest.NewRecorder()

			auh := NewAuthHandler(test.payload.cfg, logger, authRepo, auth)
			auh.Signup(rw, req)

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}

func TestSignin(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		body     []byte
		repoMock func(repoMock *repomock.MockAuth)
		authMock func(authMock *authmock.MockAuth)
	}
	type expected struct {
		code int
		body string
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Sign in with success",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"username":"test","password":"password"}`),
				repoMock: func(repoMock *repomock.MockAuth) {
					repoMock.EXPECT().GetUserByName(gomock.Any()).Return(&entity.Credentials{Password: "$2a$08$ZtefSglA0MuPtOYRa/dZI.zb.pf.dhUHo1XXmhTrKmUuMz.9Cqg6m"}, nil)
				},
				authMock: func(authMock *authmock.MockAuth) {
					authMock.EXPECT().CreateTokens(gomock.Any()).Return(&entity.Tokens{}, nil)
				},
			},
			expected: expected{
				code: http.StatusOK,
				body: `{"access_token":"","refresh_token":""}`,
			},
		},
		{
			name: "Sign in invalid body with fail",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`invalid`),
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid character 'i' looking for beginning of value"}`,
			},
		},
		{
			name: "Sign in GetUserByName error not found with fail",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"username":"test","password":"password"}`),
				repoMock: func(repoMock *repomock.MockAuth) {
					repoMock.EXPECT().GetUserByName(gomock.Any()).Return(&entity.Credentials{}, entity.ErrUserNotFound)
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
			},
			expected: expected{
				code: http.StatusNotFound,
				body: `{"error":"user not found"}`,
			},
		},
		{
			name: "Sign in GetUserByName error with fail",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"username":"test","password":"password"}`),
				repoMock: func(repoMock *repomock.MockAuth) {
					repoMock.EXPECT().GetUserByName(gomock.Any()).Return(&entity.Credentials{}, errors.New("error"))
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
				body: `{}`,
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

			authRepo := repomock.NewMockAuth(mockCtrl)
			test.payload.repoMock(authRepo)

			auth := authmock.NewMockAuth(mockCtrl)
			test.payload.authMock(auth)

			req, err := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer(test.payload.body))
			if err != nil {
				t.Error("failed to create request")
			}
			rw := httptest.NewRecorder()

			auh := NewAuthHandler(test.payload.cfg, logger, authRepo, auth)
			auh.Signin(rw, req)

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}

func TestRefresh(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		body     []byte
		repoMock func(repoMock *repomock.MockAuth)
		authMock func(authMock *authmock.MockAuth)
	}
	type expected struct {
		code int
		body string
	}
	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Refresh in with success",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"access_token":"","refresh_token":""}`),
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
					authMock.EXPECT().RefreshTokens(gomock.Any()).Return(&entity.Tokens{}, nil)
				},
			},
			expected: expected{
				code: http.StatusOK,
				body: `{"access_token":"","refresh_token":""}`,
			},
		},
		{
			name: "Refresh in invalid body with fail",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`invalid`),
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid character 'i' looking for beginning of value"}`,
			},
		},
		{
			name: "Refresh RefreshTokens error with fail",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"access_token":"","refresh_token":""}`),
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
					authMock.EXPECT().RefreshTokens(gomock.Any()).Return(&entity.Tokens{}, errors.New("error"))
				},
			},
			expected: expected{
				code: http.StatusUnauthorized,
				body: `{"error":"error"}`,
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

			authRepo := repomock.NewMockAuth(mockCtrl)
			test.payload.repoMock(authRepo)

			auth := authmock.NewMockAuth(mockCtrl)
			test.payload.authMock(auth)

			req, err := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer(test.payload.body))
			if err != nil {
				t.Error("failed to create request")
			}
			rw := httptest.NewRecorder()

			auh := NewAuthHandler(test.payload.cfg, logger, authRepo, auth)
			auh.Refresh(rw, req)

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}

func TestLogout(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		repoMock func(repoMock *repomock.MockAuth)
		authMock func(authMock *authmock.MockAuth)
		ctxMock  func(req *http.Request) context.Context
	}
	type expected struct {
		code int
		body string
	}
	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Logout with success",
			payload: payload{
				cfg: &config.Config{},
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
					authMock.EXPECT().RemoveTokens(gomock.Any(), gomock.Any()).Return(nil)
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "userUUID")
					ctx = context.WithValue(ctx, middleware.AccessUUID, "accessUUID")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusNoContent,
				body: `{}`,
			},
		},
		{
			name: "Logout invalid accessUUID with fail",
			payload: payload{
				cfg: &config.Config{},
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.Background()
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"accessUUID not found"}`,
			},
		},
		{
			name: "Logout invalid userUUID with fail",
			payload: payload{
				cfg: &config.Config{},
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.AccessUUID, "accessUUID")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"userUUID not found"}`,
			},
		},
		{
			name: "Logout RemoveTokens error with fail",
			payload: payload{
				cfg: &config.Config{},
				repoMock: func(repoMock *repomock.MockAuth) {
				},
				authMock: func(authMock *authmock.MockAuth) {
					authMock.EXPECT().RemoveTokens(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "userUUID")
					ctx = context.WithValue(ctx, middleware.AccessUUID, "accessUUID")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"error"}`,
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

			authRepo := repomock.NewMockAuth(mockCtrl)
			test.payload.repoMock(authRepo)

			auth := authmock.NewMockAuth(mockCtrl)
			test.payload.authMock(auth)

			req, err := http.NewRequest(http.MethodPost, "url", nil)
			if err != nil {
				t.Error("failed to create request")
			}

			ctx := test.payload.ctxMock(req)

			rw := httptest.NewRecorder()

			auh := NewAuthHandler(test.payload.cfg, logger, authRepo, auth)
			auh.Logout(rw, req.WithContext(ctx))

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}
