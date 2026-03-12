package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	DatabasePath   string
	JWTSecret      string
	JWTExpiryHours int
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           8080,
		DatabasePath:   "bitewise.db",
		JWTExpiryHours: 72,
	}

	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT %q: %w", v, err)
		}
		cfg.Port = p
	}

	if v := os.Getenv("DATABASE_PATH"); v != "" {
		cfg.DatabasePath = v
	}

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		h, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS %q: %w", v, err)
		}
		cfg.JWTExpiryHours = h
	}

	return cfg, nil
}
