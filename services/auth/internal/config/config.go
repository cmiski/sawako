package config

import "os"

type Config struct {
	DatabaseURL string
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

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable",
		),
	}
}
