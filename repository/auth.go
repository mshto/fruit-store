package repository

import (
	"database/sql"

	"github.com/mshto/fruit-store/backend/entity"
)

// Auth Auth
type Auth interface {
	Signup(creds *entity.Credentials) error
	GetUserByName(userName string) (*entity.Credentials, error)
}

// NewAuth NewAuth
func NewAuth(db *sql.DB) Auth {
	return &authImpl{
		db: db,
	}
}

type authImpl struct {
	db *sql.DB
}

var (
	getUserPasswordByName = "SELECT * FROM users WHERE username=$1"
	validateUserByName    = "SELECT exists (SELECT id FROM users WHERE username=$1)"
	signup                = "INSERT INTO users (username, password) VALUES ($1, $2)"
)

// Signup Signup
func (aui *authImpl) Signup(creds *entity.Credentials) error {
	exists, err := aui.isRowExist(creds.Username)
	if err != nil {
		return err
	}
	if exists {
		return entity.ErrUserAlreadyExist
	}

	_, err = aui.db.Query(signup, creds.Username, creds.Password)
	return err
}

// GetUserPasswordByName GetUserPasswordByName
func (aui *authImpl) GetUserByName(userName string) (*entity.Credentials, error) {
	var creds entity.Credentials
	err := aui.db.QueryRow(getUserPasswordByName, userName).Scan(&creds.ID, &creds.Username, &creds.Password)
	if err == sql.ErrNoRows {
		return &creds, entity.ErrUserNotFound
	}
	return &creds, err
}

func (aui *authImpl) isRowExist(username string) (bool, error) {
	var exists bool
	err := aui.db.QueryRow(validateUserByName, username).Scan(&exists)
	return exists, err
}
