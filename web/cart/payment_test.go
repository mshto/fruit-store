package cart

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	billmock "github.com/mshto/fruit-store/bill/mock"
	"github.com/mshto/fruit-store/config"
	repomock "github.com/mshto/fruit-store/repository/mock"
	"github.com/mshto/fruit-store/web/middleware"
)

func TestAddPayment(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		body     []byte
		repoMock func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount)
		billMock func(billMock *billmock.MockBill)
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
			name: "Add payment with success",
			payload: payload{
				cfg:  &config.Config{},
				body: []byte(`{"number":"number"}`),
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().RemoveUserProducts(gomock.Any()).Return(nil)
				},
				billMock: func(billMock *billmock.MockBill) {
					billMock.EXPECT().ValidateCard(gomock.Any()).Return(nil)
					billMock.EXPECT().RemoveDiscount(gomock.Any()).Return(nil)
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusCreated,
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

			cartRepo := repomock.NewMockCart(mockCtrl)
			discRepo := repomock.NewMockDiscount(mockCtrl)

			test.payload.repoMock(cartRepo, discRepo)

			billMock := billmock.NewMockBill(mockCtrl)
			test.payload.billMock(billMock)

			req, _ := http.NewRequest(http.MethodGet, "url", bytes.NewBuffer(test.payload.body))
			rw := httptest.NewRecorder()

			ctx := test.payload.ctxMock(req)

			crh := NewCardHandler(test.payload.cfg, logger, cartRepo, discRepo, billMock)
			crh.AddPayment(rw, req.WithContext(ctx))

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}
