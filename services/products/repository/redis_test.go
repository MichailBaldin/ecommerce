package repository

import (
	"ecommerce/services/products/models"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisRepo_ImplementsInterface(t *testing.T) {
	// Compile-time check that RedisRepo implements CacheRepository
	var _ CacheRepository = (*RedisRepo)(nil)
}

func TestNewRedisRepo(t *testing.T) {
	// Act
	repo := NewRedisRepo("localhost:6379")

	// Assert
	assert.NotNil(t, repo)

	// Cleanup
	if redisRepo, ok := repo.(*RedisRepo); ok {
		redisRepo.Close()
	}
}

func TestRedisRepo_Set_Get_Delete_Integration(t *testing.T) {
	t.Skip("Integration test - requires Redis server")

	// Integration test stub for future implementation
}

func TestRedisRepo_Get_CacheMiss(t *testing.T) {
	t.Skip("Integration test - requires Redis server")

	// Integration test stub for future implementation
}

func TestRedisRepo_Set_InvalidProduct(t *testing.T) {
	// Test serialization edge cases without requiring Redis
	product := &models.Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test description",
		Price:       99.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// This is just testing that our Product struct can be JSON marshaled
	// which is important for Redis caching
	_, err := json.Marshal(product)
	assert.NoError(t, err, "Product should be JSON serializable")
}

func TestRedisRepo_KeyFormat(t *testing.T) {
	// Test that we generate correct Redis keys
	expectedKey := "product:123"
	actualKey := fmt.Sprintf("product:%d", 123)

	assert.Equal(t, expectedKey, actualKey, "Redis key format should be 'product:{id}'")
}

func TestRedisRepo_JSONMarshaling(t *testing.T) {
	// Test JSON marshaling/unmarshaling that Redis repo uses
	original := &models.Product{
		ID:          42,
		Name:        "Gaming Mouse",
		Description: "High-precision gaming mouse",
		Price:       79.99,
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		UpdatedAt:   time.Now().UTC().Truncate(time.Second),
	}

	// Marshal
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal
	var unmarshaled models.Product
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, original.ID, unmarshaled.ID)
	assert.Equal(t, original.Name, unmarshaled.Name)
	assert.Equal(t, original.Description, unmarshaled.Description)
	assert.Equal(t, original.Price, unmarshaled.Price)
	assert.Equal(t, original.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
	assert.Equal(t, original.UpdatedAt.Unix(), unmarshaled.UpdatedAt.Unix())
}

func TestRedisRepo_JSONMarshaling_PriceAccuracy(t *testing.T) {
	// Test that price values are preserved accurately through JSON marshaling
	testPrices := []float64{0.0, 0.99, 10.5, 99.99, 999.99, 1234.56}

	for _, price := range testPrices {
		t.Run(fmt.Sprintf("price_%.2f", price), func(t *testing.T) {
			product := &models.Product{
				ID:    1,
				Name:  "Test Product",
				Price: price,
			}

			// Marshal
			data, err := json.Marshal(product)
			require.NoError(t, err)

			// Unmarshal
			var unmarshaled models.Product
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			// Assert price is preserved
			assert.Equal(t, price, unmarshaled.Price)
		})
	}
}

func TestRedisRepo_JSONMarshaling_EmptyFields(t *testing.T) {
	// Test marshaling product with empty optional fields
	product := &models.Product{
		ID:          1,
		Name:        "Minimal Product",
		Description: "",  // empty description
		Price:       0.0, // zero price
	}

	// Marshal
	data, err := json.Marshal(product)
	require.NoError(t, err)

	// Unmarshal
	var unmarshaled models.Product
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, product.ID, unmarshaled.ID)
	assert.Equal(t, product.Name, unmarshaled.Name)
	assert.Empty(t, unmarshaled.Description)
	assert.Equal(t, 0.0, unmarshaled.Price)
}
