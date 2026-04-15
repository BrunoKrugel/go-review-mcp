package config

import (
	"log"
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

// Config holds the server configuration including transport type,
// cache TTL, and server port.
type Config struct {
	// Transport specifies the transport type ("stdio" or "http")
	Transport string
	// CacheTTL defines how long cached content remains valid
	CacheTTL time.Duration
	// Port is the HTTP server port (used when Transport is "http")
	Port int
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	cfg := &Config{
		Transport: getEnv("TRANSPORT", "stdio"),
		Port:      getEnvInt("PORT", DefaultPort),
		CacheTTL:  getEnvDuration("CACHE_TTL", DefaultCacheTTL),
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
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("warning: invalid integer for %s=%q, using default %d: %v",
			key, value, defaultValue, err)
		return defaultValue
	}
	return i
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("warning: invalid duration for %s=%q, using default %v: %v",
			key, value, defaultValue, err)
		return defaultValue
	}
	return d
}
