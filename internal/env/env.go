package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvFile(filename string) error {
	return godotenv.Load(filename)
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
