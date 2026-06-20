package config

import "os"

type Config struct {
	Port             string
	AuthServiceURL   string
	EventServiceURL  string
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
		Port: getEnv("PORT", "8080"),
		AuthServiceURL: getEnv(
			"AUTH_SERVICE_URL",
			"http://localhost:8081",
		),
		EventServiceURL: getEnv(
			"EVENT_SERVICE_URL",
			"http://localhost:8082",
		),
	}
}
