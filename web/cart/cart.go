package cart

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/mshto/fruit-store/bill"
	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/entity"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/common/response"
	"github.com/mshto/fruit-store/web/middleware"
)

// Service Cart Service
type Service interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	AddOneProduct(w http.ResponseWriter, r *http.Request)
	RemoveProduct(w http.ResponseWriter, r *http.Request)

	AddDiscout(w http.ResponseWriter, r *http.Request)

	AddPayment(w http.ResponseWriter, r *http.Request)
}

// ProductHandler product handler struct
type cartHandler struct {
	cfg  *config.Config
	log  *logrus.Logger
	repo *repository.Repository
	bil  bill.Bill
}

// NewCardHandler NewCardHandler
func NewCardHandler(cfg *config.Config, log *logrus.Logger, repo *repository.Repository, bil bill.Bill) Service {
	return cartHandler{
		cfg:  cfg,
		log:  log,
		repo: repo,
		bil:  bil,
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

	sort.Slice(products, func(i, j int) bool {
		return products[i].Name < products[j].Name
	})

	total, err := ph.bil.GetTotalInfo(userUUID, products)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	var isDiscountAdded bool
	sale, err := ph.bil.GetDiscountByUser(userUUID)
	if err != nil && err != cache.ErrNotFound {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}
	if sale.ID != "" {
		isDiscountAdded = true
	}

	response.RenderResponse(w, http.StatusOK, entity.UserCart{
		CartProducts:    products,
		TotalPrice:      total.Price,
		TotalSavings:    total.Savings,
		Amount:          total.Amount,
		IsDiscountAdded: isDiscountAdded,
	})
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
