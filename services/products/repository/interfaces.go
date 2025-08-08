package repository

import "ecommerce/services/products/models"

// ProductRepository defines methods for product data persistence
type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id int) (*models.Product, error)
}

// CacheRepository defines methods for caching product data
type CacheRepository interface {
	Set(product *models.Product) error
	Get(id int) (*models.Product, error)
	Delete(id int) error
}
