// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ElOtro/auction-go/config"
	v1 "github.com/ElOtro/auction-go/internal/controller/http/v1"
	repo "github.com/ElOtro/auction-go/internal/infrastructure/repo/postgres"
	"github.com/ElOtro/auction-go/internal/usecase"
	"github.com/ElOtro/auction-go/pkg/httpserver"
	"github.com/ElOtro/auction-go/pkg/logger"
	"github.com/ElOtro/auction-go/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// pg models
	pgModels := repo.NewRepo(pg)

	// use cases
	useCases := usecase.NewUseCases(&pgModels)

	// controllers
	controllers := v1.NewControllers(&useCases, cfg.JWT.Secret)

	// HTTP Server
	h := v1.NewHandlers(controllers)
	httpServer := httpserver.New(h.Routes(), httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
