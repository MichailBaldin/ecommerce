package service

import (
	"ecommerce/services/users/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock implementations
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Set(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockCacheRepository) Get(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockCacheRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_ImplementsInterface(t *testing.T) {
	// Compile-time check
	var _ UserService = (*userService)(nil)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	req := &models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	mockUserRepo.On("Create", mock.MatchedBy(func(user *models.User) bool {
		return user.Name == "John Doe" && user.Email == "john@example.com"
	})).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = 123 // Simulate database setting ID
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
	}).Return(nil)

	mockCacheRepo.On("Set", mock.MatchedBy(func(user *models.User) bool {
		return user.ID == 123
	})).Return(nil)

	// Act
	result, err := service.CreateUser(req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 123, result.ID)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_DatabaseError(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	req := &models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("database error"))

	// Act
	result, err := service.CreateUser(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertNotCalled(t, "Set") // Should not try to cache on database error
}

func TestUserService_CreateUser_CacheError(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	req := &models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock expectations
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = 123
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
	}).Return(nil)

	mockCacheRepo.On("Set", mock.AnythingOfType("*models.User")).Return(errors.New("cache error"))

	// Act
	result, err := service.CreateUser(req)

	// Assert
	require.NoError(t, err) // Cache error should not fail the operation
	require.NotNil(t, result)
	assert.Equal(t, 123, result.ID)

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestUserService_GetUser_CacheHit(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	expectedUser := &models.User{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(expectedUser, nil)

	// Act
	result, err := service.GetUser(123)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)

	mockCacheRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "GetByID") // Should not hit database on cache hit
}

func TestUserService_GetUser_CacheMiss_DatabaseHit(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	expectedUser := &models.User{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(nil, nil) // Cache miss
	mockUserRepo.On("GetByID", 123).Return(expectedUser, nil)
	mockCacheRepo.On("Set", expectedUser).Return(nil) // Update cache

	// Act
	result, err := service.GetUser(123)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)

	mockCacheRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUser_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	// Mock expectations
	mockCacheRepo.On("Get", 999).Return(nil, nil)    // Cache miss
	mockUserRepo.On("GetByID", 999).Return(nil, nil) // User not found

	// Act
	result, err := service.GetUser(999)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result) // Should return nil for not found

	mockCacheRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertNotCalled(t, "Set") // Should not cache nil result
}

func TestUserService_GetUser_CacheError_DatabaseHit(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	expectedUser := &models.User{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(nil, errors.New("cache connection error"))
	mockUserRepo.On("GetByID", 123).Return(expectedUser, nil)
	mockCacheRepo.On("Set", expectedUser).Return(nil)

	// Act
	result, err := service.GetUser(123)

	// Assert
	require.NoError(t, err) // Cache error should not fail the operation
	require.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)

	mockCacheRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUser_DatabaseError(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(nil, nil) // Cache miss
	mockUserRepo.On("GetByID", 123).Return(nil, errors.New("database connection error"))

	// Act
	result, err := service.GetUser(123)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection error")

	mockCacheRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestNewUserService(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()

	// Act
	service := NewUserService(mockUserRepo, mockCacheRepo, logger)

	// Assert
	assert.NotNil(t, service)
	assert.Implements(t, (*UserService)(nil), service)
}
