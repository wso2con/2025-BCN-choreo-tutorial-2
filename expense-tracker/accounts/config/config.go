package config

import (
	"fmt"
	"os"
	"strings"
)

// Config represents the application configuration
type Config struct {
	DBType     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPath     string // For SQLite
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite
	}

	dbType = strings.ToLower(dbType)
	if dbType != "mysql" && dbType != "sqlite" {
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	config := &Config{
		DBType: dbType,
	}

	switch dbType {
	case "mysql":
		config.DBHost = getEnv("DB_HOST", "localhost")
		config.DBPort = getEnv("DB_PORT", "3306")
		config.DBUser = getEnv("DB_USER", "root")
		config.DBPassword = getEnv("DB_PASSWORD", "")
		config.DBName = getEnv("DB_NAME", "accounts")
	case "sqlite":
		config.DBPath = getEnv("DB_PATH", "./accounts.db")
	}

	return config, nil
}

// getEnv gets the environment variable or returns the fallback value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
