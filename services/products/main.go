package main

import (
	"ecommerce/services/products/config"
	"ecommerce/services/products/handlers"
	"ecommerce/services/products/repository"
	"ecommerce/services/products/service"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting products service",
		zap.String("port", cfg.Port),
		zap.String("environment", cfg.Environment))

	// Initialize repositories
	productRepo, err := repository.NewPostgresRepo(cfg.PostgresURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	cacheRepo := repository.NewRedisRepo(cfg.RedisAddr)

	// Initialize service
	productService := service.NewProductService(productRepo, cacheRepo, logger)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService, logger)

	// Setup routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/v1/products", productHandler.CreateProduct)
	http.HandleFunc("/api/v1/products/", productHandler.GetProduct)

	// Start server
	logger.Info("Server starting", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"service":"products","status":"ok"}`)
}
