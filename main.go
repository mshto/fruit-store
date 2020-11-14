package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/urfave/negroni"

	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/config"
	"github.com/mshto/fruit-store/database"
	"github.com/mshto/fruit-store/logger"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web"
	"github.com/mshto/fruit-store/web/middleware"
)

var (
	configPath      = "fruit_store_cfg.json"
	salesConfigPath = "fruit_store_sales_cfg.json"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
	}()

	config, err := config.New(configPath, salesConfigPath)
	if err != nil {
		log.Fatalf("failed to read config, error: %v", err)
	}

	log, err := logger.New(config.Logger)
	if err != nil {
		log.Fatalf("failed to setup logger, error: %v", err)
	}
	log.Infof("config: %v", config)
	db, err := database.New(config.Database)
	if err != nil {
		log.Fatalf("failed to setup db, error: %v", err)
	}

	redis, err := cache.New(config.Redis)
	if err != nil {
		log.Fatalf("failed to setup redis, error: %v", err)
	}

	repo := repository.New(db)

	router := web.New(config, log, repo, redis)
	serverMiddleware := setWebServerMiddleware()
	serverMiddleware.UseHandler(router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ListenURL),
		Handler: serverMiddleware,
	}

	log.Infof("listening for requests on port: %s", server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		e := server.Shutdown(ctx)
		log.Infof("stop running application: %v, shutdown: %v", err, e)
		os.Exit(1)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sig
	e := server.Shutdown(ctx)
	fmt.Println(e)
	log.Infof("graceful shutdown: %v, %v", err, e)
}

// move to middleware
func setWebServerMiddleware() *negroni.Negroni {
	middlewareManager := negroni.New()
	middlewareManager.Use(negroni.NewRecovery())
	middlewareManager.Use(middleware.NewWithCORSMiddleware())

	return middlewareManager
}
