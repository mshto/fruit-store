package authentication

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
)

//go:generate mockgen -destination=mock/authentication.go -package=authmock github.com/mshto/fruit-store/authentication Auth

// AccessDetails access details struct
type AccessDetails struct {
	AccessUUID string
	UserUUID   string
}

// TokenDetails token details struct
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// Auth interface
type Auth interface {
	GetUserUUID(accessUUID string) (string, error)
	CreateTokens(userUUID uuid.UUID) (*entity.Tokens, error)
	RefreshTokens(refreshToken string) (*entity.Tokens, error)
	ValidateToken(token string) (*AccessDetails, error)
	RemoveTokens(accessUUID, userUUID string) error
}

// New generate a new Auth
func New(cfg *config.Config, log *logrus.Logger, cache cache.Cache) Auth {
	return &authImpl{
		cfg:   cfg,
		log:   log,
		cache: cache,
	}
}

type authImpl struct {
	cache cache.Cache
	cfg   *config.Config
	log   *logrus.Logger
}

func (aui *authImpl) GetUserUUID(accessUUID string) (string, error) {
	return aui.cache.Get(accessUUID)
}

// CreateTokens create user tokens
func (aui *authImpl) CreateTokens(userUUID uuid.UUID) (*entity.Tokens, error) {
	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(time.Duration(aui.cfg.Auth.AccessSecretAtExpiresInMin) * time.Minute).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = td.AccessUUID + "++" + userUUID.String()

	err := aui.createAccessToken(userUUID, td)
	if err != nil {
		return nil, err
	}

	err = aui.createRefreshToken(userUUID, td)
	if err != nil {
		return nil, err
	}

	err = aui.createAuth(userUUID, td)

	return &entity.Tokens{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}, err
}

// RefreshTokens refresh tokens
func (aui *authImpl) RefreshTokens(refreshToken string) (*entity.Tokens, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(aui.cfg.Auth.RefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("refresh token is expired")
	}
	refreshUUID, ok := claims["refresh_uuid"].(string)
	if !ok {
		return nil, errors.New("refresh token is invalid")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("userUUIS is invalid")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	err = aui.cache.Del(refreshUUID)
	if err != nil {
		return nil, err
	}

	return aui.CreateTokens(userUUID)
}

// ValidateToken validate token
func (aui *authImpl) ValidateToken(tokenString string) (*AccessDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(aui.cfg.Auth.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return nil, errors.New("token is invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token is invalid")
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, err
	}
	userUUID, ok := claims["user_id"].(string)
	if !ok {
		return nil, err
	}
	return &AccessDetails{
		AccessUUID: accessUUID,
		UserUUID:   userUUID,
	}, nil
}

// RemoveTokens remove tokens
func (aui *authImpl) RemoveTokens(accessUUID, userUUID string) error {
	err := aui.cache.Del(accessUUID)
	if err != nil {
		return err
	}

	err = aui.cache.Del(fmt.Sprintf("%s++%s", accessUUID, userUUID))
	return err
}

func (aui *authImpl) createAccessToken(userUUID uuid.UUID, td *TokenDetails) error {
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userUUID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(aui.cfg.Auth.AccessSecret))
	return err
}

func (aui *authImpl) createRefreshToken(userUUID uuid.UUID, td *TokenDetails) error {
	var err error

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userUUID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(aui.cfg.Auth.RefreshSecret))
	return err
}

func (aui *authImpl) createAuth(userUUID uuid.UUID, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := aui.cache.Set(td.AccessUUID, userUUID.String(), at.Sub(now))
	if err != nil {
		return err
	}

	err = aui.cache.Set(td.RefreshUUID, userUUID.String(), rt.Sub(now))
	return err
}
