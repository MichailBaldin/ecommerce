package service

import "ecommerce/services/products/models"

type ProductService interface {
	CreateProduct(req *models.CreateProductRequest) (*models.Product, error)
	GetProduct(id int) (*models.Product, error)
}
