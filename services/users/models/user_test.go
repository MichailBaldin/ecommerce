package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_JSONMarshaling(t *testing.T) {
	// Arrange
	now := time.Now().UTC().Truncate(time.Second) // убираем наносекунды для точного сравнения
	user := User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act - Marshal to JSON
	jsonData, err := json.Marshal(user)
	require.NoError(t, err)

	// Act - Unmarshal back
	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, user.ID, unmarshaled.ID)
	assert.Equal(t, user.Name, unmarshaled.Name)
	assert.Equal(t, user.Email, unmarshaled.Email)
	assert.Equal(t, user.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
	assert.Equal(t, user.UpdatedAt.Unix(), unmarshaled.UpdatedAt.Unix())
}

func TestCreateUserRequest_JSONMarshaling(t *testing.T) {
	// Arrange
	req := CreateUserRequest{
		Name:  "Jane Smith",
		Email: "jane@example.com",
	}

	// Act
	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	var unmarshaled CreateUserRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, req.Name, unmarshaled.Name)
	assert.Equal(t, req.Email, unmarshaled.Email)
}

func TestUserResponse_JSONMarshaling(t *testing.T) {
	// Arrange
	now := time.Now().UTC().Truncate(time.Second)
	resp := UserResponse{
		ID:        42,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	jsonData, err := json.Marshal(resp)
	require.NoError(t, err)

	var unmarshaled UserResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, resp.ID, unmarshaled.ID)
	assert.Equal(t, resp.Name, unmarshaled.Name)
	assert.Equal(t, resp.Email, unmarshaled.Email)
	assert.Equal(t, resp.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
}

func TestErrorResponse_JSONMarshaling(t *testing.T) {
	// Arrange
	errResp := ErrorResponse{
		Error:   "validation_error",
		Message: "Invalid email format",
	}

	// Act
	jsonData, err := json.Marshal(errResp)
	require.NoError(t, err)

	var unmarshaled ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, errResp.Error, unmarshaled.Error)
	assert.Equal(t, errResp.Message, unmarshaled.Message)
}

func TestCreateUserRequest_EmptyFields(t *testing.T) {
	// Test empty request
	req := CreateUserRequest{}

	// Should marshal/unmarshal without errors
	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	var unmarshaled CreateUserRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Empty(t, unmarshaled.Name)
	assert.Empty(t, unmarshaled.Email)
}
