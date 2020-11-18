package bill

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/mshto/fruit-store/entity"
)

func TestGetUserProducts(t *testing.T) {
	type expected struct {
		products []entity.GetUserProduct
		isErr    bool
	}
	type payload struct {
		sqlMock func(ucMock sqlmock.Sqlmock)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		// {
		// 	name: "Get user products with success",
		// 	expected: expected{
		// 		products: []entity.GetUserProduct{getUserProductOne},
		// 		isErr:    false,
		// 	},
		// 	payload: payload{
		// 		sqlMock: func(mock sqlmock.Sqlmock) {
		// 			rows := sqlmock.NewRows([]string{"users_cart.amount", "products.id", "products.name", "products.price"}).
		// 				AddRow(getUserProductOne.Amount, getUserProductOne.ProductUUID, getUserProductOne.Name, getUserProductOne.Price)

		// 			mock.ExpectQuery("SELECT users_cart.amount, products.id, products.name, products.price FROM users_cart").WillReturnRows(rows)
		// 		},
		// 	},
		// },
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// New(cfg *config.Config, log *logrus.Logger, cache *cache.Cache)

			// db, mock, err := sqlmock.New()
			// if err != nil {
			// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			// }
			// defer db.Close()

			// test.payload.sqlMock(mock)

			// products, err := NewCartProduct(db).GetUserProducts(userUUID)
			// assert.Equal(t, products, test.expected.products)
			// if test.expected.isErr {
			// 	assert.NotNil(t, err)
			// }

			// GetTotalInfo
		})
	}
}
