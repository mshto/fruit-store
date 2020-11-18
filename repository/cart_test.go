package repository

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mshto/fruit-store/entity"
	"github.com/stretchr/testify/assert"
)

var (
	userUUID          = uuid.New()
	getUserProductOne = entity.GetUserProduct{
		ProductUUID: uuid.New(),
		Name:        "Product 1",
		Price:       10.0,
		Amount:      1,
	}
	userProductOne = entity.UserProduct{
		ProductUUID: uuid.New(),
		UserID:      userUUID,
		Amount:      1,
	}
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
		{
			name: "Get user products with success",
			expected: expected{
				products: []entity.GetUserProduct{getUserProductOne},
				isErr:    false,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"users_cart.amount", "products.id", "products.name", "products.price"}).
						AddRow(getUserProductOne.Amount, getUserProductOne.ProductUUID, getUserProductOne.Name, getUserProductOne.Price)

					mock.ExpectQuery("SELECT users_cart.amount, products.id, products.name, products.price FROM users_cart").WillReturnRows(rows)
				},
			},
		},
		{
			name: "Get user products wrong field with failed",
			expected: expected{
				products: []entity.GetUserProduct{},
				isErr:    true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "name", "wrong"}).
						AddRow(getUserProductOne.Amount, getUserProductOne.ProductUUID, ErrDB)

					mock.ExpectQuery("SELECT users_cart.amount, products.id, products.name, products.price FROM users_cart").WillReturnRows(rows)
				},
			},
		},
		{
			name: "Get user products db error with failed",
			expected: expected{
				products: []entity.GetUserProduct{},
				isErr:    true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT users_cart.amount, products.id, products.name, products.price FROM users_cart").WillReturnError(ErrNotFound)
				},
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.payload.sqlMock(mock)

			products, err := NewCartProduct(db).GetUserProducts(userUUID)
			assert.Equal(t, products, test.expected.products)
			if test.expected.isErr {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestCreateUserProducts(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		sqlMock func(ucMock sqlmock.Sqlmock)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Create user products with success",
			expected: expected{
				err: nil,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id"}).
						AddRow(userUUID)
					mock.ExpectQuery("INSERT INTO users_cart").WithArgs(userProductOne.UserID, userProductOne.ProductUUID, userProductOne.Amount).
						WillReturnRows(rows)
				},
			},
		},
		{
			name: "Create user products with failed",
			expected: expected{
				err: ErrNotFound,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("INSERT INTO users_cart").WithArgs(userProductOne.UserID, userProductOne.ProductUUID, userProductOne.Amount).
						WillReturnError(ErrNotFound)
				},
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.payload.sqlMock(mock)

			err = NewCartProduct(db).CreateUserProducts(userProductOne.UserID, userProductOne)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestCreateUserProduct(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		sqlMock func(ucMock sqlmock.Sqlmock)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Create user product with success",
			expected: expected{
				err: nil,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id"}).
						AddRow(userUUID)
					mock.ExpectQuery("INSERT INTO users_cart").WithArgs(userProductOne.UserID, userProductOne.ProductUUID, userProductOne.Amount).
						WillReturnRows(rows)
				},
			},
		},
		{
			name: "Create user product with failed",
			expected: expected{
				err: ErrNotFound,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("INSERT INTO users_cart").WithArgs(userProductOne.UserID, userProductOne.ProductUUID, userProductOne.Amount).
						WillReturnError(ErrNotFound)
				},
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.payload.sqlMock(mock)

			err = NewCartProduct(db).CreateUserProduct(userProductOne.UserID, userProductOne.ProductUUID)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestRemoveUserProducts(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		sqlMock func(ucMock sqlmock.Sqlmock)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Remove user products with success",
			expected: expected{
				err: nil,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("DELETE FROM users_cart").WithArgs(userProductOne.UserID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
		},
		{
			name: "Remove user products with failed",
			expected: expected{
				err: ErrNotFound,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("DELETE FROM users_cart").WithArgs(userProductOne.UserID).
						WillReturnError(ErrNotFound)
				},
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.payload.sqlMock(mock)

			err = NewCartProduct(db).RemoveUserProducts(userProductOne.UserID)
			assert.Equal(t, err, test.expected.err)
		})
	}
}

func TestRemoveUserProduct(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		sqlMock func(ucMock sqlmock.Sqlmock)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Remove user product with success",
			expected: expected{
				err: nil,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("DELETE FROM users_cart").WithArgs(userProductOne.UserID, userProductOne.ProductUUID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
		},
		{
			name: "Remove user product with failed",
			expected: expected{
				err: ErrNotFound,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("DELETE FROM users_cart").WithArgs(userProductOne.UserID, userProductOne.ProductUUID).
						WillReturnError(ErrNotFound)
				},
			},
		},
	}

	for _, test := range tc {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.payload.sqlMock(mock)

			err = NewCartProduct(db).RemoveUserProduct(userProductOne.UserID, userProductOne.ProductUUID)
			assert.Equal(t, err, test.expected.err)
		})
	}
}
