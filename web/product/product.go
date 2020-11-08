package product

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/common/response"
	"github.com/mshto/fruit-store/web/middleware"
)

// Service Partner Attribute Service
type Service interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	// Create(w http.ResponseWriter, r *http.Request)
	// GetOne(w http.ResponseWriter, r *http.Request)
	// Update(w http.ResponseWriter, r *http.Request)
	// Delete(w http.ResponseWriter, r *http.Request)
}

// ProductHandler ProductHandler
type productHandler struct {
	cfg  *config.Config
	log  *logrus.Logger
	repo *repository.Repository
}

// NewProductHandler NewProductHandler
func NewProductHandler(cfg *config.Config, log *logrus.Logger, repo *repository.Repository) Service {
	return productHandler{
		cfg:  cfg,
		log:  log,
		repo: repo,
	}
}

// GetAllProducts retrieves all products from db
func (ph productHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID := ctx.Value(middleware.UserUUID).(string)
	ph.log.Error(userUUID)
	products, err := ph.repo.Product.GetAll()
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
	}

	response.RenderResponse(w, http.StatusOK, products)
}
