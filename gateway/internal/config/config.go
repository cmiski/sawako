package config

import "os"

type Config struct {
	Port string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}
	return value
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
	}
}
