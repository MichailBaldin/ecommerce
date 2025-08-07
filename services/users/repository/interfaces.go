package repository

import "ecommerce/services/users/models"

// UserRepository defines methods for user data persistence
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
}

// CacheRepository defines methods for caching user data
type CacheRepository interface {
	Set(user *models.User) error
	Get(id int) (*models.User, error)
	Delete(id int) error
}
