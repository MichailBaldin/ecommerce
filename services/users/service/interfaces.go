package service

import "ecommerce/services/users/models"

type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUser(id int) (*models.User, error)
}
