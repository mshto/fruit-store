package entity

import (
	"errors"

	"github.com/google/uuid"
)

// auth errors
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user with current name is already exist")
)

// Credentials struct
type Credentials struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	PasswordRepeat string    `json:"passwordRepeat"`
}

// Tokens struct
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
