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
	configPath = "fruit_store_cfg.json"
)

func main() {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wg.Wait()
	}()

	config, err := config.New(configPath)
	if err != nil {
		log.Fatal("failed to read config, error: %w", err)
	}

	log, err := logger.New(config.Logger)
	if err != nil {
		log.Fatalf("failed to setup logger, error: %w", err)
	}

	db, err := database.New(config.Database)
	if err != nil {
		log.Fatalf("failed to setup db, error: %w", err)
	}

	redis, err := cache.New(config.Redis)
	if err != nil {
		log.Fatalf("failed to setup redis, error: %w", err)
	}

	repo := repository.New(db)

	router := web.New(config, log, repo, redis)
	serverMiddleware := setWebServerMiddleware()
	serverMiddleware.UseHandler(router)

	server := &http.Server{
		Addr:    config.ListenURL,
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
