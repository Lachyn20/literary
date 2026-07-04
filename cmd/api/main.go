package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/hemra-siirow/literary/docs"
	"github.com/hemra-siirow/literary/internal/di"
	"github.com/hemra-siirow/literary/internal/infrastructure/config"
)

// @title Literary Backend API
// @version 1.0
// @description API for Hemra Şirow literary and biography website
// @host localhost:8081
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

	listener, actualAddr, err := listenWithFallback(cfg.ServerPort)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
	defer listener.Close()

	srv := &http.Server{
		Addr:    actualAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("listening on %s", actualAddr)
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

func listenWithFallback(preferredPort string) (net.Listener, string, error) {
	if preferredPort == "" {
		preferredPort = "8080"
	}

	listener, err := net.Listen("tcp", ":"+preferredPort)
	if err == nil {
		return listener, ":" + preferredPort, nil
	}

	if !isAddressInUseError(err) {
		return nil, "", err
	}

	fallbackListener, fallbackErr := net.Listen("tcp", ":0")
	if fallbackErr != nil {
		return nil, "", fallbackErr
	}

	addr := fallbackListener.Addr().String()
	if host, port, splitErr := net.SplitHostPort(addr); splitErr == nil && host == "" {
		addr = ":" + port
	}
	return fallbackListener, addr, nil
}

func isAddressInUseError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "address already in use") || strings.Contains(msg, "only one usage of each socket address")
}
