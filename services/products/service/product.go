package service

import (
	"ecommerce/services/products/models"
	"ecommerce/services/products/repository"

	"go.uber.org/zap"
)

type productService struct {
	productRepo repository.ProductRepository
	cacheRepo   repository.CacheRepository
	logger      *zap.Logger
}

func NewProductService(productRepo repository.ProductRepository, cacheRepo repository.CacheRepository, logger *zap.Logger) ProductService {
	return &productService{
		productRepo: productRepo,
		cacheRepo:   cacheRepo,
		logger:      logger,
	}
}

func (s *productService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	s.logger.Info("Creating product",
		zap.String("name", req.Name),
		zap.Float64("price", req.Price))

	if err := s.productRepo.Create(product); err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}

	// Cache the product
	if err := s.cacheRepo.Set(product); err != nil {
		s.logger.Warn("Failed to cache product", zap.Int("id", product.ID), zap.Error(err))
		// Cache failure is not critical - continue
	}

	s.logger.Info("Product created successfully", zap.Int("id", product.ID))
	return product, nil
}

func (s *productService) GetProduct(id int) (*models.Product, error) {
	s.logger.Debug("Getting product", zap.Int("id", id))

	// Try cache first
	product, err := s.cacheRepo.Get(id)
	if err != nil {
		s.logger.Warn("Cache error", zap.Int("id", id), zap.Error(err))
		// Cache error - continue to database
	}

	if product != nil {
		s.logger.Debug("Cache hit", zap.Int("id", id))
		return product, nil
	}

	s.logger.Debug("Cache miss", zap.Int("id", id))

	// Get from database
	product, err = s.productRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get product from database", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	if product == nil {
		s.logger.Info("Product not found", zap.Int("id", id))
		return nil, nil
	}

	// Update cache
	if err := s.cacheRepo.Set(product); err != nil {
		s.logger.Warn("Failed to update cache", zap.Int("id", id), zap.Error(err))
		// Cache failure is not critical - continue
	}

	s.logger.Debug("Product retrieved from database", zap.Int("id", id))
	return product, nil
}
