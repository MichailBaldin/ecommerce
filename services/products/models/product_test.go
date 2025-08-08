package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct_JSONMarshaling(t *testing.T) {
	// Arrange
	now := time.Now().UTC().Truncate(time.Second) // убираем наносекунды для точного сравнения
	product := Product{
		ID:          1,
		Name:        "Gaming Laptop",
		Description: "High-performance gaming laptop",
		Price:       1299.99,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Act - Marshal to JSON
	jsonData, err := json.Marshal(product)
	require.NoError(t, err)

	// Act - Unmarshal back
	var unmarshaled Product
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, product.ID, unmarshaled.ID)
	assert.Equal(t, product.Name, unmarshaled.Name)
	assert.Equal(t, product.Description, unmarshaled.Description)
	assert.Equal(t, product.Price, unmarshaled.Price)
	assert.Equal(t, product.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
	assert.Equal(t, product.UpdatedAt.Unix(), unmarshaled.UpdatedAt.Unix())
}

func TestCreateProductRequest_JSONMarshaling(t *testing.T) {
	// Arrange
	req := CreateProductRequest{
		Name:        "Wireless Mouse",
		Description: "Ergonomic wireless mouse",
		Price:       29.99,
	}

	// Act
	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	var unmarshaled CreateProductRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, req.Name, unmarshaled.Name)
	assert.Equal(t, req.Description, unmarshaled.Description)
	assert.Equal(t, req.Price, unmarshaled.Price)
}

func TestProductResponse_JSONMarshaling(t *testing.T) {
	// Arrange
	now := time.Now().UTC().Truncate(time.Second)
	resp := ProductResponse{
		ID:          42,
		Name:        "Mechanical Keyboard",
		Description: "RGB mechanical keyboard",
		Price:       159.99,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Act
	jsonData, err := json.Marshal(resp)
	require.NoError(t, err)

	var unmarshaled ProductResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, resp.ID, unmarshaled.ID)
	assert.Equal(t, resp.Name, unmarshaled.Name)
	assert.Equal(t, resp.Description, unmarshaled.Description)
	assert.Equal(t, resp.Price, unmarshaled.Price)
	assert.Equal(t, resp.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
}

func TestErrorResponse_JSONMarshaling(t *testing.T) {
	// Arrange
	errResp := ErrorResponse{
		Error:   "validation_error",
		Message: "Invalid price value",
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

func TestCreateProductRequest_EmptyFields(t *testing.T) {
	// Test empty request
	req := CreateProductRequest{}

	// Should marshal/unmarshal without errors
	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	var unmarshaled CreateProductRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Empty(t, unmarshaled.Name)
	assert.Empty(t, unmarshaled.Description)
	assert.Zero(t, unmarshaled.Price)
}

func TestProduct_PriceHandling(t *testing.T) {
	// Test different price values
	testCases := []struct {
		name  string
		price float64
	}{
		{"zero price", 0.0},
		{"small price", 0.99},
		{"normal price", 99.99},
		{"large price", 9999.99},
		{"precise price", 123.456},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			product := Product{
				ID:    1,
				Name:  "Test Product",
				Price: tc.price,
			}

			jsonData, err := json.Marshal(product)
			require.NoError(t, err)

			var unmarshaled Product
			err = json.Unmarshal(jsonData, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, tc.price, unmarshaled.Price)
		})
	}
}
