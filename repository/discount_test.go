package repository

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/mshto/fruit-store/config"
	"github.com/stretchr/testify/assert"
)

var (
	saleOne = config.GeneralSale{
		ID: "codeID",
		Elements: map[string]int{
			"Oranges": 1,
		},
		Rule:     "rule",
		Discount: 10,
	}
)

func TestGetDiscount(t *testing.T) {
	type expected struct {
		sale  config.GeneralSale
		isErr bool
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
			name: "Get discount with success",
			expected: expected{
				sale:  saleOne,
				isErr: false,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"exists"}).
						AddRow(true)
					mock.ExpectQuery("SELECT exists").WithArgs(saleOne.ID).WillReturnRows(rows)

					rows = sqlmock.NewRows([]string{"id", "rule", "elements", "discount"}).
						AddRow(saleOne.ID, saleOne.Rule, []byte(`{"Oranges":1}`), saleOne.Discount)
					mock.ExpectQuery("SELECT id, rule, elements, discount FROM discount").WithArgs(saleOne.ID).WillReturnRows(rows)
				},
			},
		},
		{
			name: "Get discount row exist error with failed",
			expected: expected{
				sale:  config.GeneralSale{},
				isErr: true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT exists").WithArgs(saleOne.ID).WillReturnError(ErrNotFound)
				},
			},
		},
		{
			name: "Get discount row not exist with failed",
			expected: expected{
				sale:  config.GeneralSale{},
				isErr: true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"exists"}).
						AddRow(false)
					mock.ExpectQuery("SELECT exists").WithArgs(saleOne.ID).WillReturnRows(rows)
				},
			},
		},
		{
			name: "Get discount db err with failed",
			expected: expected{
				sale:  config.GeneralSale{},
				isErr: false,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"exists"}).
						AddRow(true)
					mock.ExpectQuery("SELECT exists").WithArgs(saleOne.ID).WillReturnRows(rows)

					rows = sqlmock.NewRows([]string{"id", "rule", "elements", "discount"}).
						AddRow(saleOne.ID, saleOne.Rule, []byte(`{"Oranges":1}`), saleOne.Discount)
					mock.ExpectQuery("SELECT id, rule, elements, discount FROM discount").WillReturnError(ErrNotFound)
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

			sale, err := NewDiscount(db).GetDiscount(saleOne.ID)
			assert.Equal(t, sale, test.expected.sale)
			if test.expected.isErr {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestRemoveDiscount(t *testing.T) {
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
			name: "Remove discount with success",
			expected: expected{
				err: nil,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("DELETE FROM discount").WithArgs(saleOne.ID).
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
		},
		{
			name: "Remove discount with failed",
			expected: expected{
				err: ErrNotFound,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectExec("DELETE FROM discount").WithArgs(saleOne.ID).
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

			err = NewDiscount(db).RemoveDiscount(saleOne.ID)
			assert.Equal(t, err, test.expected.err)
		})
	}
}
