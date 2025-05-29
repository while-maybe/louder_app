// Example in your main.go or a config package
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Common library for .env files
)

// AppConfig holds the entire configuration for the app
type AppConfig struct {
	ServerPort string
}

// LoadConfig attempt to load .env file. In production, variables are usually set directly.
func LoadConfig() *AppConfig {
	err := godotenv.Load() // Tries to load .env from the current directory or parent dirs - By default, godotenv.Load() WILL NOT OVERRIDE existing environment variables
	if err != nil {
		log.Println("no .env file found or error loading it, trying env vars")
	}

	return &AppConfig{
		ServerPort: getEnv("REST_API_SERVER_PORT", ":8080"),
	}
}

// getEnv is a helper function to get an enviroment variable or return a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
