package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/mshto/fruit-store/authentication"
	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/auth"
	"github.com/mshto/fruit-store/web/middleware"
	"github.com/mshto/fruit-store/web/product"
	// "gitlab.connectwisedev.com/platform/extended-attributes-retrieval-service/cfg"
	// "gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/logger"
)

// New creates a router for URL-to-service mapping
// func New(logger logger.Log, config *cfg.Config) *mux.Router {
func New(cfg *config.Config, log *logrus.Logger, repo *repository.Repository, redis *cache.Cache) *mux.Router {
	jwt := authentication.New(redis)

	pdh := product.NewProductHandler(cfg, log, repo)
	auh := auth.NewAuthHandler(cfg, log, repo, jwt)

	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix(cfg.URLPrefix).Subrouter()

	// healthHandler := health.NewHandler(cfg)
	// // health.NewHandler(*cfg)
	// //healthHandler := health.NewHandler(config, session)

	// api.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	// 	api := router.PathPrefix(config.URLPrefix).Subrouter()

	// 	// api.HandleFunc("/version", rest.HandlerVersion).Methods(http.MethodGet)
	// 	// api.HandleFunc("/health", rest.HandlerHealth).Methods(http.MethodGet)

	routerV1 := api.PathPrefix("/v1").Subrouter()
	routerV1.HandleFunc("/signup", auh.Signup).Methods(http.MethodPost)
	routerV1.HandleFunc("/signin", auh.Signin).Methods(http.MethodPost)
	routerV1.HandleFunc("/refresh", auh.Refresh).Methods(http.MethodPost)

	routerV1Auth := api.PathPrefix("/v1").Subrouter()
	routerV1Auth.Use(middleware.AuthMiddleware(jwt))
	// middlewares.ServeHTTP(appComponents.Log)
	// routerV1.Use(middleware.AuthMiddleware)
	routerV1Auth.HandleFunc("/logout", auh.Logout).Methods(http.MethodPost)
	routerV1Auth.HandleFunc("/products", pdh.GetAll).Methods(http.MethodGet)
	// 	// routerV1.HandleFunc(integratorsDefinitions, ias.Create).Methods(http.MethodPost)
	// 	// routerV1.HandleFunc(integratorsDefinitions+attributeID, ias.GetOne).Methods(http.MethodGet)
	// 	// routerV1.HandleFunc(integratorsDefinitions+attributeID, ias.Update).Methods(http.MethodPut)
	// 	// routerV1.HandleFunc(integratorsDefinitions+attributeID, ias.Delete).Methods(http.MethodDelete)

	return router
}
