package cart

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	loggermock "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/mshto/fruit-store/bill"
	billmock "github.com/mshto/fruit-store/bill/mock"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
	repomock "github.com/mshto/fruit-store/repository/mock"
	"github.com/mshto/fruit-store/web/middleware"
)

func TestGetAll(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		url      string
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
			name: "Get all with success",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().GetUserProducts(gomock.Any()).Return([]entity.GetUserProduct{
						{Name: "Second"}, {Name: "First"},
					}, nil)
				},
				billMock: func(billMock *billmock.MockBill) {
					billMock.EXPECT().GetTotalInfo(gomock.Any(), gomock.Any()).Return(bill.TotalInfo{}, nil)
					billMock.EXPECT().GetDiscountByUser(gomock.Any()).Return(config.GeneralSale{ID: "sale ID"}, nil)
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusOK,
				body: `{"products":[{"id":"00000000-0000-0000-0000-000000000000","name":"First","price":0,"amount":0},{"id":"00000000-0000-0000-0000-000000000000","name":"Second","price":0,"amount":0}],"totalPrice":"","totalSavings":"","totalAmount":"","isDiscountAdded":true}`,
			},
		},
		{
			name: "Get all invalid user uuid with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid UUID length: 8"}`},
		},
		{
			name: "Get all GetUserProducts error with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().GetUserProducts(gomock.Any()).Return([]entity.GetUserProduct{}, errors.New("error"))
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
				body: `{"error":"error"}`,
			},
		},
		{
			name: "Get all GetTotalInfo error with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().GetUserProducts(gomock.Any()).Return([]entity.GetUserProduct{}, nil)
				},
				billMock: func(billMock *billmock.MockBill) {
					billMock.EXPECT().GetTotalInfo(gomock.Any(), gomock.Any()).Return(bill.TotalInfo{}, errors.New("error"))
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusInternalServerError,
				body: `{"error":"error"}`,
			},
		},
		{
			name: "Get all GetDiscountByUser error with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().GetUserProducts(gomock.Any()).Return([]entity.GetUserProduct{}, nil)
				},
				billMock: func(billMock *billmock.MockBill) {
					billMock.EXPECT().GetTotalInfo(gomock.Any(), gomock.Any()).Return(bill.TotalInfo{}, nil)
					billMock.EXPECT().GetDiscountByUser(gomock.Any()).Return(config.GeneralSale{}, errors.New("error"))
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
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

			cartRepo := repomock.NewMockCart(mockCtrl)
			discRepo := repomock.NewMockDiscount(mockCtrl)

			test.payload.repoMock(cartRepo, discRepo)

			billMock := billmock.NewMockBill(mockCtrl)
			test.payload.billMock(billMock)

			req, _ := http.NewRequest(http.MethodGet, test.payload.url, nil)
			rw := httptest.NewRecorder()

			ctx := test.payload.ctxMock(req)

			crh := NewCardHandler(test.payload.cfg, logger, cartRepo, discRepo, billMock)

			router := mux.NewRouter()
			router.HandleFunc("/v1/cart/products", crh.GetAll)
			router.ServeHTTP(rw, req.WithContext(ctx))

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		url      string
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
			name: "Update product with success",
			payload: payload{
				cfg:  &config.Config{},
				url:  "/v1/cart/products",
				body: []byte(`{}`),
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().CreateUserProducts(gomock.Any(), gomock.Any()).Return(nil)
				},
				billMock: func(billMock *billmock.MockBill) {
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
		{
			name: "Update product invalid user uuid with fail",
			payload: payload{
				cfg:  &config.Config{},
				url:  "/v1/cart/products",
				body: []byte(`{}`),
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid UUID length: 8"}`,
			},
		},
		{
			name: "Update product invalid body with fail",
			payload: payload{
				cfg:  &config.Config{},
				url:  "/v1/cart/products",
				body: []byte(`invalid`),
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid character 'i' looking for beginning of value"}`,
			},
		},
		{
			name: "Update product db error with fail",
			payload: payload{
				cfg:  &config.Config{},
				url:  "/v1/cart/products",
				body: []byte(`{}`),
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().CreateUserProducts(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
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

			cartRepo := repomock.NewMockCart(mockCtrl)
			discRepo := repomock.NewMockDiscount(mockCtrl)

			test.payload.repoMock(cartRepo, discRepo)

			billMock := billmock.NewMockBill(mockCtrl)
			test.payload.billMock(billMock)

			req, _ := http.NewRequest(http.MethodGet, test.payload.url, bytes.NewBuffer(test.payload.body))
			rw := httptest.NewRecorder()

			ctx := test.payload.ctxMock(req)

			crh := NewCardHandler(test.payload.cfg, logger, cartRepo, discRepo, billMock)

			router := mux.NewRouter()
			router.HandleFunc("/v1/cart/products", crh.UpdateProduct)
			router.ServeHTTP(rw, req.WithContext(ctx))

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}

func TestAddOneProduct(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		url      string
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
			name: "Add one product with success",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480-2c1a-11eb-adc1-0242ac120002",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().CreateUserProduct(gomock.Any(), gomock.Any()).Return(nil)
				},
				billMock: func(billMock *billmock.MockBill) {
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
		{
			name: "Add one product invalid user UUID with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480-2c1a-11eb-adc1-0242ac120002",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid UUID length: 8"}`,
			},
		},
		{
			name: "Add one product invalid product id with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid UUID length: 8"}`,
			},
		},
		{
			name: "Add one product db errorwith fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480-2c1a-11eb-adc1-0242ac120002",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().CreateUserProduct(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
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

			cartRepo := repomock.NewMockCart(mockCtrl)
			discRepo := repomock.NewMockDiscount(mockCtrl)

			test.payload.repoMock(cartRepo, discRepo)

			billMock := billmock.NewMockBill(mockCtrl)
			test.payload.billMock(billMock)

			req, _ := http.NewRequest(http.MethodGet, test.payload.url, nil)
			rw := httptest.NewRecorder()

			ctx := test.payload.ctxMock(req)

			crh := NewCardHandler(test.payload.cfg, logger, cartRepo, discRepo, billMock)

			router := mux.NewRouter()
			router.HandleFunc("/v1/cart/products/{productID}", crh.AddOneProduct)
			router.ServeHTTP(rw, req.WithContext(ctx))

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}

func TestRemoveProduct(t *testing.T) {
	type payload struct {
		cfg      *config.Config
		url      string
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
			name: "Remove product with success",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480-2c1a-11eb-adc1-0242ac120002",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().RemoveUserProduct(gomock.Any(), gomock.Any()).Return(nil)
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusNoContent,
				body: `{}`,
			},
		},
		{
			name: "Remove product invalid user UUID with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480-2c1a-11eb-adc1-0242ac120002",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid UUID length: 8"}`,
			},
		},
		{
			name: "Remove product invalid product id with fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
				},
			},
			expected: expected{
				code: http.StatusBadRequest,
				body: `{"error":"invalid UUID length: 8"}`,
			},
		},
		{
			name: "Remove product db errorwith fail",
			payload: payload{
				cfg: &config.Config{},
				url: "/v1/cart/products/e2d49480-2c1a-11eb-adc1-0242ac120002",
				repoMock: func(cartMock *repomock.MockCart, discMock *repomock.MockDiscount) {
					cartMock.EXPECT().RemoveUserProduct(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				},
				billMock: func(billMock *billmock.MockBill) {
				},
				ctxMock: func(req *http.Request) context.Context {
					ctx := context.WithValue(req.Context(), middleware.UserUUID, "e2d49480-2c1a-11eb-adc1-0242ac120002")
					return ctx
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

			cartRepo := repomock.NewMockCart(mockCtrl)
			discRepo := repomock.NewMockDiscount(mockCtrl)

			test.payload.repoMock(cartRepo, discRepo)

			billMock := billmock.NewMockBill(mockCtrl)
			test.payload.billMock(billMock)

			req, _ := http.NewRequest(http.MethodGet, test.payload.url, nil)
			rw := httptest.NewRecorder()

			ctx := test.payload.ctxMock(req)

			crh := NewCardHandler(test.payload.cfg, logger, cartRepo, discRepo, billMock)

			router := mux.NewRouter()
			router.HandleFunc("/v1/cart/products/{productID}", crh.RemoveProduct)
			router.ServeHTTP(rw, req.WithContext(ctx))

			assert.Equal(t, test.expected.code, rw.Code)
			assert.Equal(t, test.expected.body, rw.Body.String())
		})
	}
}
