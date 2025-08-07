package repository

import (
	"ecommerce/services/users/models"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_InterfaceContract(t *testing.T) {
	// Test that UserRepository interface has expected methods
	userRepoType := reflect.TypeOf((*UserRepository)(nil)).Elem()

	// Check interface has correct number of methods
	assert.Equal(t, 2, userRepoType.NumMethod(), "UserRepository should have exactly 2 methods")

	// Check method names exist
	methods := make(map[string]bool)
	for i := 0; i < userRepoType.NumMethod(); i++ {
		method := userRepoType.Method(i)
		methods[method.Name] = true
	}

	assert.True(t, methods["Create"], "UserRepository should have Create method")
	assert.True(t, methods["GetByID"], "UserRepository should have GetByID method")
}

func TestCacheRepository_InterfaceContract(t *testing.T) {
	// Test that CacheRepository interface has expected methods
	cacheRepoType := reflect.TypeOf((*CacheRepository)(nil)).Elem()

	// Check interface has correct number of methods
	assert.Equal(t, 3, cacheRepoType.NumMethod(), "CacheRepository should have exactly 3 methods")

	// Check method names exist
	methods := make(map[string]bool)
	for i := 0; i < cacheRepoType.NumMethod(); i++ {
		method := cacheRepoType.Method(i)
		methods[method.Name] = true
	}

	assert.True(t, methods["Set"], "CacheRepository should have Set method")
	assert.True(t, methods["Get"], "CacheRepository should have Get method")
	assert.True(t, methods["Delete"], "CacheRepository should have Delete method")
}

// Mock implementations for compile-time interface checks
type mockUserRepository struct{}

func (m *mockUserRepository) Create(user *models.User) error       { return nil }
func (m *mockUserRepository) GetByID(id int) (*models.User, error) { return nil, nil }

type mockCacheRepository struct{}

func (m *mockCacheRepository) Set(user *models.User) error      { return nil }
func (m *mockCacheRepository) Get(id int) (*models.User, error) { return nil, nil }
func (m *mockCacheRepository) Delete(id int) error              { return nil }

func TestMockImplementations_CompileTimeCheck(t *testing.T) {
	// These assignments will fail to compile if interfaces are not implemented correctly
	var userRepo UserRepository = &mockUserRepository{}
	var cacheRepo CacheRepository = &mockCacheRepository{}

	assert.NotNil(t, userRepo, "mockUserRepository should implement UserRepository")
	assert.NotNil(t, cacheRepo, "mockCacheRepository should implement CacheRepository")
}

func TestInterfaceMethodSignatures(t *testing.T) {
	// Test specific method signatures for UserRepository
	userRepoType := reflect.TypeOf((*UserRepository)(nil)).Elem()

	// Check Create method signature
	createMethod, exists := userRepoType.MethodByName("Create")
	assert.True(t, exists, "Create method should exist")
	assert.Equal(t, 1, createMethod.Type.NumIn(), "Create should take 1 parameter (user)")
	assert.Equal(t, 1, createMethod.Type.NumOut(), "Create should return 1 value (error)")

	// Check GetByID method signature
	getMethod, exists := userRepoType.MethodByName("GetByID")
	assert.True(t, exists, "GetByID method should exist")
	assert.Equal(t, 1, getMethod.Type.NumIn(), "GetByID should take 1 parameter (id)")
	assert.Equal(t, 2, getMethod.Type.NumOut(), "GetByID should return 2 values (*User, error)")
}
