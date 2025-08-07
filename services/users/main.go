package main

import (
	"log"
	"net/http"

	"ecommerce/services/users/config"
	"ecommerce/services/users/handlers"
	"ecommerce/services/users/repository"
	"ecommerce/services/users/service"

	"go.uber.org/zap"
)

func main() {
	// Config
	cfg := config.New()

	// Logger
	var logger *zap.Logger
	var err error

	if cfg.Environment == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	logger.Info("Starting user-service",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment))

	// Repositories
	userRepo, err := repository.NewPostgresRepo(cfg.PostgresURL)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	cacheRepo := repository.NewRedisRepo(cfg.RedisAddr)

	// Services
	userService := service.NewUserService(userRepo, cacheRepo, logger)

	// Handlers
	userHandler := handlers.NewUserHandler(userService, logger)

	// Routes
	http.HandleFunc("/api/v1/users", userHandler.CreateUser)
	http.HandleFunc("/api/v1/users/", userHandler.GetUser)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok","service":"users"}`)); err != nil {
			logger.Error("Failed to write health response", zap.Error(err))
		}
	})

	logger.Info("User service started successfully", zap.String("port", cfg.Port))

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
