package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/hemra-siirow/literary/internal/di"
	"github.com/hemra-siirow/literary/internal/infrastructure/config"
	_ "github.com/hemra-siirow/literary/docs"
)

// @title Literary Backend API
// @version 1.0
// @description API for Hemra Şirow literary and biography website
// @host localhost:8080
// @basePath /
// @schemes http https

func main() {
	cfg := config.Load()
	container, err := di.NewContainer(cfg)
	if err != nil {
		log.Fatalf("di error: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", container.Router)
	// swagger UI — requires running `swag init` to generate docs
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: mux,
	}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	if err := container.Shutdown(ctx); err != nil {
		log.Printf("container shutdown err: %v", err)
	}
}
