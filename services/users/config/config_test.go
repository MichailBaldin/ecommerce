package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_NewWithDefaults(t *testing.T) {
	// Arrange - очищаем environment variables
	clearTestEnvVars()

	// Act
	cfg := New()

	// Assert - проверяем значения по умолчанию
	assert.Equal(t, "8001", cfg.Port)
	assert.Equal(t, "users", cfg.ServiceName)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "development", cfg.Environment)
	assert.Contains(t, cfg.PostgresURL, "host=localhost")
	assert.Contains(t, cfg.PostgresURL, "port=5432")
	assert.Contains(t, cfg.PostgresURL, "user=postgres")
	assert.Contains(t, cfg.PostgresURL, "dbname=users")
	assert.Contains(t, cfg.PostgresURL, "sslmode=disable")
}

func TestConfig_NewWithEnvironmentVariables(t *testing.T) {
	// Arrange - устанавливаем тестовые environment variables
	setupTestEnvVars()
	defer clearTestEnvVars()

	// Act
	cfg := New()

	// Assert - проверяем что переменные окружения применились
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "test-users", cfg.ServiceName)
	assert.Equal(t, "redis:6380", cfg.RedisAddr)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "testing", cfg.Environment)
	assert.Contains(t, cfg.PostgresURL, "host=testhost")
	assert.Contains(t, cfg.PostgresURL, "port=5433")
	assert.Contains(t, cfg.PostgresURL, "user=testuser")
	assert.Contains(t, cfg.PostgresURL, "password=testpass")
	assert.Contains(t, cfg.PostgresURL, "dbname=testdb")
}

func TestBuildPostgresURL_AllDefaults(t *testing.T) {
	// Arrange - очищаем переменные БД
	clearPostgresEnvVars()

	// Act
	url := buildPostgresURL()

	// Assert
	expected := "host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable"
	assert.Equal(t, expected, url)
}

func TestBuildPostgresURL_WithCustomValues(t *testing.T) {
	// Arrange - устанавливаем кастомные значения
	os.Setenv("POSTGRES_HOST", "custom-host")
	os.Setenv("POSTGRES_PORT", "5434")
	os.Setenv("POSTGRES_USER", "custom-user")
	os.Setenv("POSTGRES_PASSWORD", "custom-pass")
	os.Setenv("POSTGRES_DB", "custom-db")
	defer clearPostgresEnvVars()

	// Act
	url := buildPostgresURL()

	// Assert
	expected := "host=custom-host port=5434 user=custom-user password=custom-pass dbname=custom-db sslmode=disable"
	assert.Equal(t, expected, url)
}

func TestGetEnv_WithValue(t *testing.T) {
	// Arrange
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	// Act
	result := getEnv("TEST_VAR", "default")

	// Assert
	assert.Equal(t, "test-value", result)
}

func TestGetEnv_WithoutValue(t *testing.T) {
	// Arrange - убеждаемся что переменной нет
	os.Unsetenv("NON_EXISTENT_VAR")

	// Act
	result := getEnv("NON_EXISTENT_VAR", "default-value")

	// Assert
	assert.Equal(t, "default-value", result)
}

func TestGetEnv_WithEmptyValue(t *testing.T) {
	// Arrange - устанавливаем пустое значение
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	// Act
	result := getEnv("EMPTY_VAR", "default")

	// Assert
	assert.Equal(t, "default", result) // пустая строка должна использовать default
}

// Helper functions
func setupTestEnvVars() {
	os.Setenv("PORT", "8080")
	os.Setenv("SERVICE_NAME", "test-users")
	os.Setenv("REDIS_ADDR", "redis:6380")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENVIRONMENT", "testing")
	os.Setenv("POSTGRES_HOST", "testhost")
	os.Setenv("POSTGRES_PORT", "5433")
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("POSTGRES_DB", "testdb")
}

func clearTestEnvVars() {
	vars := []string{
		"PORT", "SERVICE_NAME", "REDIS_ADDR", "LOG_LEVEL", "ENVIRONMENT",
	}
	for _, v := range vars {
		os.Unsetenv(v)
	}
	clearPostgresEnvVars()
}

func clearPostgresEnvVars() {
	postgresVars := []string{
		"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
	}
	for _, v := range postgresVars {
		os.Unsetenv(v)
	}
}
