package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/controller"
	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/logger"
)

func main() {
	container := di.NewContainer()

	loggerInst := logger.ProvideLogger(container)
	serverConfig := config.ProvideServerConfig(container)

	loggerInst.Info("Starting application on " + serverConfig.ServerEndpoint)

	router := controller.ProvideRouter(container)

	serverChan := make(chan error, 1)

	srv := &http.Server{
		Addr: serverConfig.ServerEndpoint,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "//")
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			router.ServeHTTP(w, r)
		}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverChan <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		loggerInst.Info("Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			loggerInst.Error(fmt.Sprintf("Server shutdown error: %v", err))
		}
	case err := <-serverChan:
		loggerInst.Error(fmt.Sprintf("Server error: %v", err))
	}

}
