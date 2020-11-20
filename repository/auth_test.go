package repository

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/mshto/fruit-store/entity"
	"github.com/stretchr/testify/assert"
)

var (
	cred = entity.Credentials{
		ID:       userUUID,
		Username: "test",
		Password: "password",
	}
)

func TestGetUserByName(t *testing.T) {
	type expected struct {
		cred  entity.Credentials
		isErr bool
	}
	type payload struct {
		sqlMock func(sqlMock sqlmock.Sqlmock)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "GetUserByName with success",
			expected: expected{
				cred:  cred,
				isErr: false,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "username", "password"}).
						AddRow(cred.ID, cred.Username, cred.Password)
					mock.ExpectQuery("SELECT id, username, password FROM users").WithArgs(cred.Username).WillReturnRows(rows)
				},
			},
		},
		{
			name: "GetUserByName db error with failed",
			expected: expected{
				cred:  entity.Credentials{},
				isErr: true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT id, username, password FROM users").WithArgs(cred.Username).WillReturnError(ErrNotFound)
				},
			},
		},
		{
			name: "GetUserByName ErrNoRows error with failed",
			expected: expected{
				cred:  entity.Credentials{},
				isErr: true,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT id, username, password FROM users").WithArgs(cred.Username).WillReturnError(sql.ErrNoRows)
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

			creds, err := NewAuth(db).GetUserByName(cred.Username)
			assert.Equal(t, *creds, test.expected.cred)
			if test.expected.isErr {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestSignup(t *testing.T) {
	type expected struct {
		err error
	}
	type payload struct {
		sqlMock func(sqlMock sqlmock.Sqlmock)
	}

	tc := []struct {
		name string
		expected
		payload
	}{
		{
			name: "Signup with success",
			expected: expected{
				err: nil,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"exists"}).
						AddRow(false)
					mock.ExpectQuery("SELECT exists").WithArgs(cred.Username).WillReturnRows(rows)

					rows = sqlmock.NewRows([]string{"id"}).
						AddRow(userUUID)
					mock.ExpectQuery("INSERT INTO users").WithArgs(cred.Username, cred.Password).
						WillReturnRows(rows)
				},
			},
		},
		{
			name: "Signup db error with failed",
			expected: expected{
				err: ErrNotFound,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT exists").WithArgs(cred.Username).WillReturnError(ErrNotFound)
				},
			},
		},
		{
			name: "Signup already exist with failed",
			expected: expected{
				err: entity.ErrUserAlreadyExist,
			},
			payload: payload{
				sqlMock: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"exists"}).
						AddRow(true)
					mock.ExpectQuery("SELECT exists").WithArgs(cred.Username).WillReturnRows(rows)
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

			err = NewAuth(db).Signup(&cred)
			assert.Equal(t, err, test.expected.err)
		})
	}
}
