package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// .env.local > .env
func LoadEnvFiles() error {
	_ = godotenv.Load(".env.local")

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}
	return nil
}

func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvRequired(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s is required", key)
}
