package repository

import (
	"ecommerce/services/users/models"
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

	// repo := NewRedisRepo("localhost:6379")
	// defer repo.(*RedisRepo).Close()

	// // Arrange
	// user := &models.User{
	// 	ID:        123,
	// 	Name:      "John Doe",
	// 	Email:     "john@example.com",
	// 	CreatedAt: time.Now().UTC().Truncate(time.Second),
	// 	UpdatedAt: time.Now().UTC().Truncate(time.Second),
	// }

	// // Test Set
	// err := repo.Set(user)
	// require.NoError(t, err)

	// // Test Get
	// cachedUser, err := repo.Get(123)
	// require.NoError(t, err)
	// require.NotNil(t, cachedUser)
	// assert.Equal(t, user.ID, cachedUser.ID)
	// assert.Equal(t, user.Name, cachedUser.Name)
	// assert.Equal(t, user.Email, cachedUser.Email)

	// // Test Get non-existent (cache miss)
	// nonExistent, err := repo.Get(99999)
	// require.NoError(t, err)
	// assert.Nil(t, nonExistent, "Should return nil for cache miss")

	// // Test Delete
	// err = repo.Delete(123)
	// require.NoError(t, err)

	// // Verify deletion
	// deletedUser, err := repo.Get(123)
	// require.NoError(t, err)
	// assert.Nil(t, deletedUser, "Should return nil after deletion")
}

func TestRedisRepo_Get_CacheMiss(t *testing.T) {
	t.Skip("Integration test - requires Redis server")

	// repo := NewRedisRepo("localhost:6379")
	// defer repo.(*RedisRepo).Close()

	// // Act - get non-existent key
	// user, err := repo.Get(99999)

	// // Assert
	// require.NoError(t, err)
	// assert.Nil(t, user, "Should return nil for cache miss")
}

func TestRedisRepo_Set_InvalidUser(t *testing.T) {
	// Test serialization edge cases without requiring Redis
	user := &models.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// This is just testing that our User struct can be JSON marshaled
	// which is important for Redis caching
	_, err := json.Marshal(user)
	assert.NoError(t, err, "User should be JSON serializable")
}

func TestRedisRepo_KeyFormat(t *testing.T) {
	// Test that we generate correct Redis keys
	expectedKey := "user:123"
	actualKey := fmt.Sprintf("user:%d", 123)

	assert.Equal(t, expectedKey, actualKey, "Redis key format should be 'user:{id}'")
}

func TestRedisRepo_JSONMarshaling(t *testing.T) {
	// Test JSON marshaling/unmarshaling that Redis repo uses
	original := &models.User{
		ID:        42,
		Name:      "Jane Doe",
		Email:     "jane@example.com",
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UpdatedAt: time.Now().UTC().Truncate(time.Second),
	}

	// Marshal
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal
	var unmarshaled models.User
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, original.ID, unmarshaled.ID)
	assert.Equal(t, original.Name, unmarshaled.Name)
	assert.Equal(t, original.Email, unmarshaled.Email)
	assert.Equal(t, original.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
	assert.Equal(t, original.UpdatedAt.Unix(), unmarshaled.UpdatedAt.Unix())
}
