package config

import (
	"os"
	"strconv"
	"time"
)

const (
	// DefaultPort is the default HTTP port for the server
	DefaultPort = 8080
	// DefaultCacheTTL is the default cache TTL duration
	DefaultCacheTTL = 24 * time.Hour
)

// Config holds the server configuration
type Config struct {
	Transport string
	LogLevel  string
	CacheTTL  time.Duration
	Port      int
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	cfg := &Config{
		Transport: getEnv("TRANSPORT", "stdio"),
		Port:      getEnvInt("PORT", DefaultPort),
		CacheTTL:  getEnvDuration("CACHE_TTL", DefaultCacheTTL),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
