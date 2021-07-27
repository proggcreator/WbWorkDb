package config

import (
	"os"
)

type Config struct {
	Username string
	Password string
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
func New() *Config {
	return &Config{

		Username: getEnv("Username", ""),
		Password: getEnv("Password", ""),
	}
}

// Simple helper function to read an environment or return a default value
