package repository

import (
	"ecommerce/services/products/models"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductRepository_InterfaceContract(t *testing.T) {
	// Test that ProductRepository interface has expected methods
	productRepoType := reflect.TypeOf((*ProductRepository)(nil)).Elem()

	// Check interface has correct number of methods
	assert.Equal(t, 2, productRepoType.NumMethod(), "ProductRepository should have exactly 2 methods")

	// Check method names exist
	methods := make(map[string]bool)
	for i := 0; i < productRepoType.NumMethod(); i++ {
		method := productRepoType.Method(i)
		methods[method.Name] = true
	}

	assert.True(t, methods["Create"], "ProductRepository should have Create method")
	assert.True(t, methods["GetByID"], "ProductRepository should have GetByID method")
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
type mockProductRepository struct{}

func (m *mockProductRepository) Create(product *models.Product) error    { return nil }
func (m *mockProductRepository) GetByID(id int) (*models.Product, error) { return nil, nil }

type mockCacheRepository struct{}

func (m *mockCacheRepository) Set(product *models.Product) error   { return nil }
func (m *mockCacheRepository) Get(id int) (*models.Product, error) { return nil, nil }
func (m *mockCacheRepository) Delete(id int) error                 { return nil }

func TestMockImplementations_CompileTimeCheck(t *testing.T) {
	// These assignments will fail to compile if interfaces are not implemented correctly
	var productRepo ProductRepository = &mockProductRepository{}
	var cacheRepo CacheRepository = &mockCacheRepository{}

	assert.NotNil(t, productRepo, "mockProductRepository should implement ProductRepository")
	assert.NotNil(t, cacheRepo, "mockCacheRepository should implement CacheRepository")
}

func TestInterfaceMethodSignatures(t *testing.T) {
	// Test specific method signatures for ProductRepository
	productRepoType := reflect.TypeOf((*ProductRepository)(nil)).Elem()

	// Check Create method signature
	createMethod, exists := productRepoType.MethodByName("Create")
	assert.True(t, exists, "Create method should exist")
	assert.Equal(t, 1, createMethod.Type.NumIn(), "Create should take 1 parameter (product)")
	assert.Equal(t, 1, createMethod.Type.NumOut(), "Create should return 1 value (error)")

	// Check GetByID method signature
	getMethod, exists := productRepoType.MethodByName("GetByID")
	assert.True(t, exists, "GetByID method should exist")
	assert.Equal(t, 1, getMethod.Type.NumIn(), "GetByID should take 1 parameter (id)")
	assert.Equal(t, 2, getMethod.Type.NumOut(), "GetByID should return 2 values (*Product, error)")
}
