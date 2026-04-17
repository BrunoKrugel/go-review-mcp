package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("TRANSPORT")
	os.Unsetenv("PORT")
	os.Unsetenv("CACHE_TTL")

	cfg := LoadFromEnv()

	if cfg.Transport != "stdio" {
		t.Errorf("expected transport 'stdio', got %q", cfg.Transport)
	}
	if cfg.Port != DefaultPort {
		t.Errorf("expected port %d, got %d", DefaultPort, cfg.Port)
	}
	if cfg.CacheTTL != DefaultCacheTTL {
		t.Errorf("expected cache TTL %v, got %v", DefaultCacheTTL, cfg.CacheTTL)
	}
}

func TestLoadFromEnv_Custom(t *testing.T) {
	t.Setenv("TRANSPORT", "http")
	t.Setenv("PORT", "9090")
	t.Setenv("CACHE_TTL", "12h")

	cfg := LoadFromEnv()

	if cfg.Transport != "http" {
		t.Errorf("expected transport 'http', got %q", cfg.Transport)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.CacheTTL != 12*time.Hour {
		t.Errorf("expected cache TTL 12h, got %v", cfg.CacheTTL)
	}
}

func TestLoadFromEnv_InvalidPort(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	cfg := LoadFromEnv()

	if cfg.Port != DefaultPort {
		t.Errorf("expected default port %d for invalid value, got %d", DefaultPort, cfg.Port)
	}
}

func TestLoadFromEnv_InvalidDuration(t *testing.T) {
	t.Setenv("CACHE_TTL", "not-a-duration")

	cfg := LoadFromEnv()

	if cfg.CacheTTL != DefaultCacheTTL {
		t.Errorf("expected default TTL %v for invalid value, got %v", DefaultCacheTTL, cfg.CacheTTL)
	}
}
