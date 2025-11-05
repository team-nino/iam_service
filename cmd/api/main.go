package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/team-nino/iam_service/internal/infrastructure/config"
	httpHandler "github.com/team-nino/iam_service/internal/infrastructure/http"
	"github.com/team-nino/iam_service/internal/infrastructure/persistence"
	"github.com/team-nino/iam_service/internal/usecase"
	"github.com/team-nino/iam_service/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.NewSimpleLogger()
	log.Info("Starting IAM Service...")

	// Load configuration
	cfg := config.Load()
	log.Info("Configuration loaded successfully")

	// Initialize repositories (using in-memory for now)
	userRepo := persistence.NewInMemoryUserRepository()
	sessionRepo := persistence.NewInMemorySessionRepository()
	log.Info("Repositories initialized")

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, sessionRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	log.Info("Use cases initialized")

	// Initialize handlers
	authHandler := httpHandler.NewAuthHandler(authUseCase)
	userHandler := httpHandler.NewUserHandler(userUseCase)
	log.Info("Handlers initialized")

	// Setup router
	router := httpHandler.NewRouter(authHandler, userHandler)
	mux := router.SetupRoutes()
	log.Info("Routes configured")

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
