package config

import "os"

type Config struct {
	Port        string
	ServiceName string
	PostgresURL string
	RedisAddr   string
	LogLevel    string
	Environment string
}

func New() *Config {
	return &Config{
		Port:        getEnv("PORT", "8001"),
		ServiceName: getEnv("SERVICE_NAME", "users"),
		PostgresURL: buildPostgresURL(),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func buildPostgresURL() string {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "password")
	dbname := getEnv("POSTGRES_DB", "users")

	return "host=" + host + " port=" + port + " user=" + user +
		" password=" + password + " dbname=" + dbname + " sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
