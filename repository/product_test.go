package repository

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mshto/fruit-store/entity"
	"github.com/stretchr/testify/assert"
)

var (
	ErrDB = "record not found"
)

var (
	productOne = entity.Product{
		ID:        uuid.New(),
		Name:      "Product 1",
		Price:     10.0,
		CreatedAt: time.Now(),
	}
)

func TestGetAll(t *testing.T) {
	type expected struct {
		products []entity.Product
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
			name: "Get all with success",
			expected: expected{
				products: []entity.Product{productOne},
				isErr:    false,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "name", "price", "created_at"}).
						AddRow(productOne.ID, productOne.Name, productOne.Price, productOne.CreatedAt)

					mock.ExpectQuery(getAllProducts).WillReturnRows(rows)
				},
			},
		},
		{
			name: "Get all wrong field with failed",
			expected: expected{
				products: []entity.Product{},
				isErr:    true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "name", "wrong"}).
						AddRow(productOne.ID, productOne.Name, ErrDB)

					mock.ExpectQuery(getAllProducts).WillReturnRows(rows)
				},
			},
		},
		{
			name: "Get all db error with failed",
			expected: expected{
				products: []entity.Product{},
				isErr:    true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery(getAllProducts).WillReturnError(ErrNotFound)
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

			products, err := NewProduct(db).GetAll()
			assert.Equal(t, products, test.expected.products)
			if test.expected.isErr {
				assert.NotNil(t, err)
			}
		})
	}
}
