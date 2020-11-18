package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/mshto/fruit-store/authentication"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/common/response"
	"github.com/mshto/fruit-store/web/middleware"
)

// Service Service
type Service interface {
	Signup(w http.ResponseWriter, r *http.Request)
	Signin(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	// Create(w http.ResponseWriter, r *http.Request)
	// GetOne(w http.ResponseWriter, r *http.Request)
	// Update(w http.ResponseWriter, r *http.Request)
	// Delete(w http.ResponseWriter, r *http.Request)
}

// ProductHandler ProductHandler
type authHandler struct {
	cfg  *config.Config
	log  *logrus.Logger
	repo *repository.Repository
	auth authentication.Auth
}

// NewAuthHandler NewAuthHandler
func NewAuthHandler(cfg *config.Config, log *logrus.Logger, repo *repository.Repository, auth authentication.Auth) Service {
	return authHandler{
		cfg:  cfg,
		log:  log,
		repo: repo,
		auth: auth,
	}
}

// Signup Signup
func (ah authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	creds := &entity.Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	creds.Password = string(hashedPassword)
	err = ah.repo.Auth.Signup(creds)
	// if errors.Is(err, entity.ErrUserAlreadyExist) {
	if err == entity.ErrUserAlreadyExist {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.RenderResponse(w, http.StatusCreated, response.EmptyResp{})
}

// Signin Signin
func (ah authHandler) Signin(w http.ResponseWriter, r *http.Request) {
	creds := &entity.Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	if creds.Password != creds.PasswordRepeat {
		response.RenderFailedResponse(w, http.StatusNotFound, errors.New("passwords aren't equal"))
		return
	}

	storedUser, err := ah.repo.Auth.GetUserByName(creds.Username)
	// if errors.Is(err, entity.ErrUserNotFound) {
	if err == entity.ErrUserNotFound {
		response.RenderFailedResponse(w, http.StatusNotFound, err)
		return
	}
	ah.log.Error(storedUser)
	if err != nil {
		response.RenderResponse(w, http.StatusInternalServerError, err)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		response.RenderResponse(w, http.StatusUnauthorized, response.EmptyResp{})
		return
	}

	tokens, err := ah.auth.CreateTokens(storedUser.ID)
	if err != nil {
		response.RenderResponse(w, http.StatusForbidden, err)
		return
	}

	response.RenderResponse(w, http.StatusOK, tokens)
}

// Logout Logout
func (ah authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	tokens := &entity.Tokens{}
	err := json.NewDecoder(r.Body).Decode(tokens)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}
	generatedTokens, err := ah.auth.RefreshTokens(tokens.RefreshToken)
	if err != nil {
		// If the two passwords don't match, return a 401 status
		response.RenderFailedResponse(w, http.StatusUnauthorized, err)
		return
	}
	response.RenderResponse(w, http.StatusOK, generatedTokens)
}

// Logout Logout
func (ah authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accessUUID, ok := ctx.Value(middleware.AccessUUID).(string)
	if !ok {
		response.RenderFailedResponse(w, http.StatusBadRequest, errors.New("accessUUID not found"))
		return
	}
	userUUID, ok := ctx.Value(middleware.UserUUID).(string)
	if !ok {
		response.RenderFailedResponse(w, http.StatusBadRequest, errors.New("userUUID not found"))
		return
	}
	err := ah.auth.RemoveTokens(accessUUID, userUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}
	response.RenderResponse(w, http.StatusNoContent, response.EmptyResp{})
}
