package cart

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/common/response"
	"github.com/mshto/fruit-store/web/middleware"
)

// Service Partner Attribute Service
type Service interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	AddOneProduct(w http.ResponseWriter, r *http.Request)
	RemoveProduct(w http.ResponseWriter, r *http.Request)
	// Create(w http.ResponseWriter, r *http.Request)
	// GetOne(w http.ResponseWriter, r *http.Request)
	// Update(w http.ResponseWriter, r *http.Request)
	// Delete(w http.ResponseWriter, r *http.Request)
}

// ProductHandler ProductHandler
type cartHandler struct {
	cfg  *config.Config
	log  *logrus.Logger
	repo *repository.Repository
}

// NewCardHandler NewCardHandler
func NewCardHandler(cfg *config.Config, log *logrus.Logger, repo *repository.Repository) Service {
	return cartHandler{
		cfg:  cfg,
		log:  log,
		repo: repo,
	}
}

// GetAllProducts retrieves all products from db
func (ph cartHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID, err := uuid.Parse(ctx.Value(middleware.UserUUID).(string))
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	products, err := ph.repo.Cart.GetUserProducts(userUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	total := entity.UserCart{
		CartProducts: products,
		Total:        ph.calculateTotal(products),
	}
	response.RenderResponse(w, http.StatusOK, total)
}

// GetAllProducts retrieves all products from db
func (ph cartHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID, err := uuid.Parse(ctx.Value(middleware.UserUUID).(string))
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	prd := &entity.UserProduct{}
	err = json.NewDecoder(r.Body).Decode(prd)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	err = ph.repo.Cart.CreateUserProducts(userUUID, *prd)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.RenderResponse(w, http.StatusCreated, response.EmptyResp{})
}

// GetAllProducts retrieves all products from db
func (ph cartHandler) AddOneProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID, err := uuid.Parse(ctx.Value(middleware.UserUUID).(string))
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	productUUID, err := uuid.Parse(mux.Vars(r)["productID"])
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	err = ph.repo.Cart.CreateUserProduct(userUUID, productUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.RenderResponse(w, http.StatusCreated, response.EmptyResp{})
}

// GetAllProducts retrieves all products from db
func (ph cartHandler) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID, err := uuid.Parse(ctx.Value(middleware.UserUUID).(string))
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	productUUID, err := uuid.Parse(mux.Vars(r)["productID"])
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	err = ph.repo.Cart.RemoveUserProduct(userUUID, productUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.RenderResponse(w, http.StatusNoContent, response.EmptyResp{})
}

func (ph cartHandler) calculateTotal(products []entity.GetUserProduct) string {
	var total float32
	for _, prd := range products {
		total += float32(prd.Amount) * prd.Price
	}

	return fmt.Sprintf("%.2f", total)
}
