// Example in your main.go or a config package
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv" // Common library for .env files
)

// AppConfig holds the entire configuration for the app
type AppConfig struct {
	ServerPort            string
	GeoAPIBaseURL         string
	GeoAPIKeyHeaderName   string
	GeoAPIKey             string
	GeoAPIRateLimitSleep  time.Duration
	GeoAPIPageLimit       int
	GeoAPICountryEndpoint string
}

// LoadConfig attempt to load .env file. In production, variables are usually set directly.
func LoadConfig() *AppConfig {
	err := godotenv.Load() // Tries to load .env from the current directory or parent dirs - By default, godotenv.Load() WILL NOT OVERRIDE existing environment variables
	if err != nil {
		log.Println("no .env file found or error loading it, trying env vars")
	}

	// ignore parsing error as this is just to load from .env
	parsedGeoAPIRateLimitSleep, _ := time.ParseDuration(getEnv("GEO_API_RATE_LIMIT_SLEEP", "1500ms"))

	// ignore parsing error as this is just to load from .env
	parsedGeoAPIRateLimit, _ := strconv.Atoi((getEnv("GEO_API_PAGE_LIMIT", "10")))

	return &AppConfig{
		ServerPort:            getEnv("REST_API_SERVER_PORT", "8080"),
		GeoAPIBaseURL:         getEnv("GEO_API_BASEURL", "https://wft-geo-db.p.rapidapi.com"),
		GeoAPIKeyHeaderName:   getEnv("GEO_API_KEY_HEADER_NAME", "x-rapidapi-key"),
		GeoAPIKey:             getEnv("GEO_API_KEY", "COULD_READ_GET_API_KEY"),
		GeoAPIRateLimitSleep:  parsedGeoAPIRateLimitSleep,
		GeoAPIPageLimit:       parsedGeoAPIRateLimit,
		GeoAPICountryEndpoint: getEnv("GEO_API_COUNTRY_ENDPOINT", "/v1/geo/countries"),
	}
}

// getEnv is a helper function to get an enviroment variable or return a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
