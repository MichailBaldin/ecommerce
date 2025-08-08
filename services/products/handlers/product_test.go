package handlers

import (
	"bytes"
	"ecommerce/services/products/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) GetProduct(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	reqBody := models.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedProduct := &models.Product{
		ID:          123,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("CreateProduct", mock.MatchedBy(func(req *models.CreateProductRequest) bool {
		return req.Name == "Test Product" && req.Description == "Test Description" &&
			req.Price == 99.99
	})).Return(expectedProduct, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.Product
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, response.ID)
	assert.Equal(t, expectedProduct.Name, response.Name)
	assert.Equal(t, expectedProduct.Description, response.Description)
	assert.Equal(t, expectedProduct.Price, response.Price)

	mockService.AssertExpectations(t)
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Bad Request", response.Error)
	assert.Equal(t, "Invalid JSON", response.Message)

	mockService.AssertNotCalled(t, "CreateProduct")
}

func TestProductHandler_CreateProduct_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	reqBody := models.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockService.On("CreateProduct", mock.AnythingOfType("*models.CreateProductRequest")).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	// Act
	handler.CreateProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Internal Server Error", response.Error)
	assert.Equal(t, "Failed to create product", response.Message)

	mockService.AssertExpectations(t)
}

func TestProductHandler_CreateProduct_WrongMethod(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	// Act
	handler.CreateProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Method Not Allowed", response.Error)
	assert.Equal(t, "Method not allowed", response.Message)
}

func TestProductHandler_GetProduct_Success(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	expectedProduct := &models.Product{
		ID:          123,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("GetProduct", 123).Return(expectedProduct, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.Product
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, response.ID)
	assert.Equal(t, expectedProduct.Name, response.Name)
	assert.Equal(t, expectedProduct.Description, response.Description)
	assert.Equal(t, expectedProduct.Price, response.Price)

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	mockService.On("GetProduct", 999).Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/999", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Not Found", response.Error)
	assert.Equal(t, "Product not found", response.Message)

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct_InvalidID(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/invalid", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Bad Request", response.Error)
	assert.Equal(t, "Invalid product ID", response.Message)

	mockService.AssertNotCalled(t, "GetProduct")
}

func TestProductHandler_GetProduct_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	mockService.On("GetProduct", 123).Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Internal Server Error", response.Error)
	assert.Equal(t, "Failed to get product", response.Message)

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct_WrongMethod(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()
	handler := NewProductHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetProduct(w, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestProductHandler_ExtractID(t *testing.T) {
	// Arrange
	handler := &ProductHandler{}

	tests := []struct {
		name        string
		path        string
		expectedID  int
		expectError bool
	}{
		{"valid ID", "/api/v1/products/123", 123, false},
		{"single digit", "/api/v1/products/5", 5, false},
		{"large ID", "/api/v1/products/999999", 999999, false},
		{"invalid path", "/products", 0, true},
		{"non-numeric ID", "/api/v1/products/abc", 0, true},
		{"empty ID", "/api/v1/products/", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			id, err := handler.extractID(req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestNewProductHandler(t *testing.T) {
	// Arrange
	mockService := &MockProductService{}
	logger := zap.NewNop()

	// Act
	handler := NewProductHandler(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
	assert.Equal(t, logger, handler.logger)
}
