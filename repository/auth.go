package repository

import (
	"database/sql"

	"github.com/mshto/fruit-store/entity"
)

// Auth interface
type Auth interface {
	GetUserByName(userName string) (*entity.Credentials, error)
	Signup(creds *entity.Credentials) error
}

// NewAuth generate new auth
func NewAuth(db *sql.DB) Auth {
	return &authImpl{
		db: db,
	}
}

type authImpl struct {
	db *sql.DB
}

var (
	getUserPasswordByName = "SELECT id, username, password FROM users WHERE username=$1"
	validateUserByName    = "SELECT exists (SELECT id FROM users WHERE username=$1)"
	signup                = "INSERT INTO users (username, password) VALUES ($1, $2)"
)

// GetUserByName get user creds by name
func (aui *authImpl) GetUserByName(userName string) (*entity.Credentials, error) {
	var creds entity.Credentials
	err := aui.db.QueryRow(getUserPasswordByName, userName).Scan(&creds.ID, &creds.Username, &creds.Password)
	if err == sql.ErrNoRows {
		return &creds, entity.ErrUserNotFound
	}
	return &creds, err
}

// Signup sign up user
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

func (aui *authImpl) isRowExist(username string) (bool, error) {
	var exists bool
	err := aui.db.QueryRow(validateUserByName, username).Scan(&exists)
	return exists, err
}
