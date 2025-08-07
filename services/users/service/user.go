package service

import (
	"ecommerce/services/users/models"
	"ecommerce/services/users/repository"

	"go.uber.org/zap"
)

type userService struct {
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
	logger    *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, cacheRepo repository.CacheRepository, logger *zap.Logger) UserService {
	return &userService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		logger:    logger,
	}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	s.logger.Info("Creating user", zap.String("email", req.Email))

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	// Cache the user
	if err := s.cacheRepo.Set(user); err != nil {
		s.logger.Warn("Failed to cache user", zap.Int("id", user.ID), zap.Error(err))
		// Cache failure is not critical - continue
	}

	s.logger.Info("User created successfully", zap.Int("id", user.ID))
	return user, nil
}

func (s *userService) GetUser(id int) (*models.User, error) {
	s.logger.Debug("Getting user", zap.Int("id", id))

	// Try cache first
	user, err := s.cacheRepo.Get(id)
	if err != nil {
		s.logger.Warn("Cache error", zap.Int("id", id), zap.Error(err))
		// Cache error - continue to database
	}

	if user != nil {
		s.logger.Debug("Cache hit", zap.Int("id", id))
		return user, nil
	}

	s.logger.Debug("Cache miss", zap.Int("id", id))

	// Get from database
	user, err = s.userRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get user from database", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	if user == nil {
		s.logger.Info("User not found", zap.Int("id", id))
		return nil, nil
	}

	// Update cache
	if err := s.cacheRepo.Set(user); err != nil {
		s.logger.Warn("Failed to update cache", zap.Int("id", id), zap.Error(err))
		// Cache failure is not critical - continue
	}

	s.logger.Debug("User retrieved from database", zap.Int("id", id))
	return user, nil
}
