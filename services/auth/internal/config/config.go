package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultPort                  = "8081"
	defaultAccessTokenTTLSeconds = 900
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func getEnv(
	key,
	defaultValue string,
) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvDurationSeconds(
	key string,
	defaultSeconds int,
) time.Duration {
	value := os.Getenv(key)

	if value == "" {
		return time.Duration(defaultSeconds) * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return time.Duration(defaultSeconds) * time.Second
	}

	return time.Duration(seconds) * time.Second
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", defaultPort),
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable",
		),
		JWTSecret: getEnv(
			"JWT_SECRET",
			"dev-secret-change-me",
		),
		AccessTokenTTL: getEnvDurationSeconds(
			"ACCESS_TOKEN_TTL_SECONDS",
			defaultAccessTokenTTLSeconds,
		),
		RefreshTokenTTL: getEnvDurationSeconds(
			"REFRESH_TOKEN_TTL_SECONDS",
			int((30 * 24 * time.Hour).Seconds()),
		),
	}
}
