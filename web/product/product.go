package product

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/common/response"
)

// Service product interface
type Service interface {
	GetAll(w http.ResponseWriter, r *http.Request)
}

// ProductHandler product handler
type productHandler struct {
	cfg         *config.Config
	log         *logrus.Logger
	productRepo repository.Products
}

// NewProductHandler init a new product handler
func NewProductHandler(cfg *config.Config, log *logrus.Logger, productRepo repository.Products) Service {
	return productHandler{
		cfg:         cfg,
		log:         log,
		productRepo: productRepo,
	}
}

// GetAllProducts retrieves all products
func (ph productHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := ph.productRepo.GetAll()
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.RenderResponse(w, http.StatusOK, products)
}
