package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"Prime_Number_Tester/internal/config"
	v1 "Prime_Number_Tester/internal/delivery/http/v1"
	"Prime_Number_Tester/internal/service"
)

type App struct {
	httpServer *http.Server
	numbers    service.NumbersServiceProvider
}

func NewApp(cfg *config.Config) (*App, error) {
	return &App{numbers: service.NewNumbersService(cfg.PrimesFilePath)}, nil
}

func (a *App) Run(port string) error {
	router := mux.NewRouter()
	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	v1.RegisterHTTPEndpoints(router, a.numbers)

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Failed to listen and serve", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt) //nolint
	<-quit
	log.Println("Shutting numbers service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Shutting numbers server...")
	return a.httpServer.Shutdown(ctx)
}
