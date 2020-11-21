package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/mshto/fruit-store/authentication"
	"github.com/mshto/fruit-store/bill"
	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/auth"
	"github.com/mshto/fruit-store/web/cart"
	"github.com/mshto/fruit-store/web/middleware"
	"github.com/mshto/fruit-store/web/product"
)

// New creates a router for URL-to-service mapping
func New(cfg *config.Config, log *logrus.Logger, repo *repository.Repository, redis cache.Cache) *mux.Router {
	jwt := authentication.New(cfg, log, redis)
	bil := bill.New(cfg, log, redis)

	pdh := product.NewProductHandler(cfg, log, repo.Product)
	cth := cart.NewCardHandler(cfg, log, repo.Cart, repo.Discount, bil)
	auh := auth.NewAuthHandler(cfg, log, repo.Auth, jwt)

	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix(cfg.URLPrefix).Subrouter()

	routerV1 := api.PathPrefix("/v1").Subrouter()
	routerV1.HandleFunc("/signup", auh.Signup).Methods(http.MethodPost)
	routerV1.HandleFunc("/signin", auh.Signin).Methods(http.MethodPost)
	routerV1.HandleFunc("/refresh", auh.Refresh).Methods(http.MethodPost)

	routerV1Auth := api.PathPrefix("/v1").Subrouter()
	routerV1Auth.Use(middleware.AuthMiddleware(jwt))

	routerV1Auth.HandleFunc("/logout", auh.Logout).Methods(http.MethodPost)
	routerV1Auth.HandleFunc("/products", pdh.GetAll).Methods(http.MethodGet)

	routerV1Auth.HandleFunc("/cart/products", cth.GetAll).Methods(http.MethodGet)
	routerV1Auth.HandleFunc("/cart/products", cth.UpdateProduct).Methods(http.MethodPost)
	routerV1Auth.HandleFunc("/cart/products/{productID}", cth.AddOneProduct).Methods(http.MethodPost)
	routerV1Auth.HandleFunc("/cart/products/{productID}", cth.RemoveProduct).Methods(http.MethodDelete)

	routerV1Auth.HandleFunc("/cart/discount", cth.AddDiscout).Methods(http.MethodPost)

	routerV1Auth.HandleFunc("/cart/payment", cth.AddPayment).Methods(http.MethodPost)

	return router
}
