package repository

import (
	"ecommerce/services/users/models"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_Create(t *testing.T) {
	t.Skip("Integration test - requires PostgreSQL database")

	// repo := setupTestDB(t)
	// defer repo.(*PostgresRepo).Close()

	// // Arrange
	// user := &models.User{
	// 	Name:  "John Doe",
	// 	Email: "john@example.com",
	// }

	// // Act
	// err := repo.Create(user)

	// // Assert
	// require.NoError(t, err)
	// assert.Greater(t, user.ID, 0, "ID should be set after create")
	// assert.NotZero(t, user.CreatedAt, "CreatedAt should be set")
	// assert.NotZero(t, user.UpdatedAt, "UpdatedAt should be set")
	// assert.Equal(t, user.CreatedAt.Unix(), user.UpdatedAt.Unix(), "CreatedAt and UpdatedAt should be equal on create")
}

func TestPostgresRepo_GetByID(t *testing.T) {
	t.Skip("Integration test - requires PostgreSQL database")

	// repo := setupTestDB(t)
	// defer repo.(*PostgresRepo).Close()

	// // Arrange - create a user first
	// originalUser := &models.User{
	// 	Name:  "Jane Smith",
	// 	Email: "jane@example.com",
	// }
	// err := repo.Create(originalUser)
	// require.NoError(t, err)

	// // Act
	// retrievedUser, err := repo.GetByID(originalUser.ID)

	// // Assert
	// require.NoError(t, err)
	// require.NotNil(t, retrievedUser)
	// assert.Equal(t, originalUser.ID, retrievedUser.ID)
	// assert.Equal(t, originalUser.Name, retrievedUser.Name)
	// assert.Equal(t, originalUser.Email, retrievedUser.Email)
	// assert.Equal(t, originalUser.CreatedAt.Unix(), retrievedUser.CreatedAt.Unix())
	// assert.Equal(t, originalUser.UpdatedAt.Unix(), retrievedUser.UpdatedAt.Unix())
}

func TestPostgresRepo_GetByID_NotFound(t *testing.T) {
	t.Skip("Integration test - requires PostgreSQL database")

	// repo := setupTestDB(t)
	// defer repo.(*PostgresRepo).Close()

	// // Act
	// user, err := repo.GetByID(99999) // non-existent ID

	// // Assert
	// require.NoError(t, err)
	// assert.Nil(t, user, "Should return nil for non-existent user")
}

// Unit tests that don't require database
func TestNewPostgresRepo_InvalidDSN(t *testing.T) {
	// Act
	repo, err := NewPostgresRepo("invalid-dsn")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, repo)
}

func TestPostgresRepo_ImplementsInterface(t *testing.T) {
	// Compile-time check that PostgresRepo implements UserRepository
	var _ UserRepository = (*PostgresRepo)(nil)
}

func TestUser_RequiredFields(t *testing.T) {
	// Test that User model has required fields for database operations
	user := &models.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Email)
	assert.Zero(t, user.ID)        // Should be zero before database save
	assert.Zero(t, user.CreatedAt) // Should be zero before database save
	assert.Zero(t, user.UpdatedAt) // Should be zero before database save
}
