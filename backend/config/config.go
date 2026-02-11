package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DBPath       string
	LogPath      string
	MaxLimit     int
	DefaultLimit int
	CORSEnabled  bool
	CORSOrigins  string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DBPath:       getEnv("DB_PATH", "./metrics.db"),
		LogPath:      getEnv("LOG_PATH", "./logs/app.log"),
		MaxLimit:     getEnvInt("MAX_LIMIT", 1000),
		DefaultLimit: getEnvInt("DEFAULT_LIMIT", 100),
		CORSEnabled:  getEnvBool("CORS_ENABLED", true),
		CORSOrigins:  getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}
	if c.DBPath == "" {
		return fmt.Errorf("DB_PATH cannot be empty")
	}
	if c.MaxLimit <= 0 {
		return fmt.Errorf("MAX_LIMIT must be positive")
	}
	if c.DefaultLimit <= 0 || c.DefaultLimit > c.MaxLimit {
		return fmt.Errorf("DEFAULT_LIMIT must be between 1 and MAX_LIMIT")
	}
	return nil
}
