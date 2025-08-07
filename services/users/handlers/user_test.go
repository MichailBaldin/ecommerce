package handlers

import (
	"bytes"
	"ecommerce/services/users/models"
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

// Mock UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUser(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	reqBody := models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedUser := &models.User{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("CreateUser", mock.MatchedBy(func(req *models.CreateUserRequest) bool {
		return req.Name == "John Doe" && req.Email == "john@example.com"
	})).Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.User
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Name, response.Name)
	assert.Equal(t, expectedUser.Email, response.Email)

	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Bad Request", response.Error)
	assert.Equal(t, "Invalid JSON", response.Message)

	mockService.AssertNotCalled(t, "CreateUser")
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	reqBody := models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockService.On("CreateUser", mock.AnythingOfType("*models.CreateUserRequest")).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	// Act
	handler.CreateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Internal Server Error", response.Error)
	assert.Equal(t, "Failed to create user", response.Message)

	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_WrongMethod(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()

	// Act
	handler.CreateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Method Not Allowed", response.Error)
	assert.Equal(t, "Method not allowed", response.Message)
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	expectedUser := &models.User{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetUser", 123).Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.User
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Name, response.Name)
	assert.Equal(t, expectedUser.Email, response.Email)

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	mockService.On("GetUser", 999).Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Not Found", response.Error)
	assert.Equal(t, "User not found", response.Message)

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_InvalidID(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Bad Request", response.Error)
	assert.Equal(t, "Invalid user ID", response.Message)

	mockService.AssertNotCalled(t, "GetUser")
}

func TestUserHandler_GetUser_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	mockService.On("GetUser", 123).Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Internal Server Error", response.Error)
	assert.Equal(t, "Failed to get user", response.Message)

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_WrongMethod(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/123", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestUserHandler_ExtractID(t *testing.T) {
	// Arrange
	handler := &UserHandler{}

	tests := []struct {
		name        string
		path        string
		expectedID  int
		expectError bool
	}{
		{"valid ID", "/api/v1/users/123", 123, false},
		{"single digit", "/api/v1/users/5", 5, false},
		{"large ID", "/api/v1/users/999999", 999999, false},
		{"invalid path", "/users", 0, true},
		{"non-numeric ID", "/api/v1/users/abc", 0, true},
		{"empty ID", "/api/v1/users/", 0, true},
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

func TestNewUserHandler(t *testing.T) {
	// Arrange
	mockService := &MockUserService{}
	logger := zap.NewNop()

	// Act
	handler := NewUserHandler(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
	assert.Equal(t, logger, handler.logger)
}
