package service

import (
	"ecommerce/services/products/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock implementations
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Set(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockCacheRepository) Get(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockCacheRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProductService_ImplementsInterface(t *testing.T) {
	// Compile-time check
	var _ ProductService = (*productService)(nil)
}

func TestProductService_CreateProduct_Success(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	req := &models.CreateProductRequest{
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		Price:       1299.99,
	}

	// Mock expectations
	mockProductRepo.On("Create", mock.MatchedBy(func(product *models.Product) bool {
		return product.Name == "Gaming Laptop" &&
			product.Description == "High-performance gaming laptop" &&
			product.Price == 1299.99
	})).Run(func(args mock.Arguments) {
		product := args.Get(0).(*models.Product)
		product.ID = 123 // Simulate database setting ID
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
	}).Return(nil)

	mockCacheRepo.On("Set", mock.MatchedBy(func(product *models.Product) bool {
		return product.ID == 123
	})).Return(nil)

	// Act
	result, err := service.CreateProduct(req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 123, result.ID)
	assert.Equal(t, "Gaming Laptop", result.Name)
	assert.Equal(t, "High-performance gaming laptop", result.Description)
	assert.Equal(t, 1299.99, result.Price)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)

	mockProductRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_DatabaseError(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	req := &models.CreateProductRequest{
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		Price:       1299.99,
	}

	// Mock expectations
	mockProductRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(errors.New("database error"))

	// Act
	result, err := service.CreateProduct(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockProductRepo.AssertExpectations(t)
	mockCacheRepo.AssertNotCalled(t, "Set") // Should not try to cache on database error
}

func TestProductService_CreateProduct_CacheError(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	req := &models.CreateProductRequest{
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		Price:       1299.99,
	}

	// Mock expectations
	mockProductRepo.On("Create", mock.AnythingOfType("*models.Product")).Run(func(args mock.Arguments) {
		product := args.Get(0).(*models.Product)
		product.ID = 123
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
	}).Return(nil)

	mockCacheRepo.On("Set", mock.AnythingOfType("*models.Product")).Return(errors.New("cache error"))

	// Act
	result, err := service.CreateProduct(req)

	// Assert
	require.NoError(t, err) // Cache error should not fail the operation
	require.NotNil(t, result)
	assert.Equal(t, 123, result.ID)

	mockProductRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_CacheHit(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	expectedProduct := &models.Product{
		ID:          123,
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		Price:       1299.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(expectedProduct, nil)

	// Act
	result, err := service.GetProduct(123)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedProduct.ID, result.ID)
	assert.Equal(t, expectedProduct.Name, result.Name)
	assert.Equal(t, expectedProduct.Description, result.Description)
	assert.Equal(t, expectedProduct.Price, result.Price)

	mockCacheRepo.AssertExpectations(t)
	mockProductRepo.AssertNotCalled(t, "GetByID") // Should not hit database on cache hit
}

func TestProductService_GetProduct_CacheMiss_DatabaseHit(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	expectedProduct := &models.Product{
		ID:          123,
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		Price:       1299.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(nil, nil) // Cache miss
	mockProductRepo.On("GetByID", 123).Return(expectedProduct, nil)
	mockCacheRepo.On("Set", expectedProduct).Return(nil) // Update cache

	// Act
	result, err := service.GetProduct(123)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedProduct.ID, result.ID)

	mockCacheRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_ProductNotFound(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	// Mock expectations
	mockCacheRepo.On("Get", 999).Return(nil, nil)       // Cache miss
	mockProductRepo.On("GetByID", 999).Return(nil, nil) // Product not found

	// Act
	result, err := service.GetProduct(999)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result) // Should return nil for not found

	mockCacheRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockCacheRepo.AssertNotCalled(t, "Set") // Should not cache nil result
}

func TestProductService_GetProduct_CacheError_DatabaseHit(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	expectedProduct := &models.Product{
		ID:          123,
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		Price:       1299.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(nil, errors.New("cache connection error"))
	mockProductRepo.On("GetByID", 123).Return(expectedProduct, nil)
	mockCacheRepo.On("Set", expectedProduct).Return(nil)

	// Act
	result, err := service.GetProduct(123)

	// Assert
	require.NoError(t, err) // Cache error should not fail the operation
	require.NotNil(t, result)
	assert.Equal(t, expectedProduct.ID, result.ID)

	mockCacheRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_DatabaseError(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	// Mock expectations
	mockCacheRepo.On("Get", 123).Return(nil, nil) // Cache miss
	mockProductRepo.On("GetByID", 123).Return(nil, errors.New("database connection error"))

	// Act
	result, err := service.GetProduct(123)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection error")

	mockCacheRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestNewProductService(t *testing.T) {
	// Arrange
	mockProductRepo := &MockProductRepository{}
	mockCacheRepo := &MockCacheRepository{}
	logger := zap.NewNop()

	// Act
	service := NewProductService(mockProductRepo, mockCacheRepo, logger)

	// Assert
	assert.NotNil(t, service)
	assert.Implements(t, (*ProductService)(nil), service)
}
