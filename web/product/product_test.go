package product

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
	repomock "github.com/mshto/fruit-store/repository/mock"
)

func TestGetAll(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		repoMock func(repoMock *repomock.MockProducts)
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
			name: "Get total info user sale with success",
			payload: payload{
				cfg: &config.Config{},
				repoMock: func(repoMock *repomock.MockProducts) {
					repoMock.EXPECT().GetAll().Return([]entity.Product{}, nil)
				},
			},
			expected: expected{
				code: http.StatusOK,
				body: `[]`,
			},
		},
		{
			name: "Get total info user sale with fail",
			payload: payload{
				cfg: &config.Config{},
				repoMock: func(repoMock *repomock.MockProducts) {
					repoMock.EXPECT().GetAll().Return([]entity.Product{}, errors.New("error"))
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
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

			productRepo := repomock.NewMockProducts(mockCtrl)

			test.payload.repoMock(productRepo)

			req, _ := http.NewRequest(http.MethodGet, "url", nil)
			rw := httptest.NewRecorder()

			pdh := NewProductHandler(test.payload.cfg, logger, productRepo)
			pdh.GetAll(rw, req)

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}
