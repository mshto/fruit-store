package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	loggermock "github.com/sirupsen/logrus/hooks/test"

	"github.com/mshto/fruit-store/authentication"
	authmock "github.com/mshto/fruit-store/authentication/mock"
)

func TestAuthMiddleware(t *testing.T) {
	type payload struct {
		authMock            func(mock *authmock.MockAuth)
		authorizationHeader string
	}

	tc := []struct {
		name string
		payload
	}{
		{
			name: "Auth middleware with success",
			payload: payload{
				authMock: func(mock *authmock.MockAuth) {
					mock.EXPECT().ValidateToken(gomock.Any()).Return(&authentication.AccessDetails{UserUUID: "test"}, nil)
					mock.EXPECT().GetUserUUID(gomock.Any()).Return("test", nil)
				},
				authorizationHeader: "Bearer eyJhbGciOiJIUzI1",
			},
		},
		{
			name: "Auth middleware with fail",
			payload: payload{
				authMock: func(mock *authmock.MockAuth) {
				},
				authorizationHeader: "Bearer",
			},
		},
		{
			name: "Auth middleware ValidateToken with fail",
			payload: payload{
				authMock: func(mock *authmock.MockAuth) {
					mock.EXPECT().ValidateToken(gomock.Any()).Return(nil, errors.New("error"))
				},
				authorizationHeader: "Bearer eyJhbGciOiJIUzI1",
			},
		},
		{
			name: "Auth middleware GetUserUUID with fail",
			payload: payload{
				authMock: func(mock *authmock.MockAuth) {
					mock.EXPECT().ValidateToken(gomock.Any()).Return(&authentication.AccessDetails{UserUUID: "test"}, nil)
					mock.EXPECT().GetUserUUID(gomock.Any()).Return("", errors.New("error"))
				},
				authorizationHeader: "Bearer eyJhbGciOiJIUzI1",
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			auth := authmock.NewMockAuth(mockCtrl)
			test.payload.authMock(auth)

			logger, _ := loggermock.NewNullLogger()

			req, _ := http.NewRequest(http.MethodGet, "url", nil)
			req.Header.Set("Authorization", test.payload.authorizationHeader)

			rw := httptest.NewRecorder()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			AuthMiddleware(auth, logger)(next).ServeHTTP(rw, req)
		})
	}
}
