// Example in your main.go or a config package
package config

import (
	"log"
	"os"
)

func main() {
	// ...
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://default.api.com/items" // Default value
		log.Println("Warning: API_URL not set, using default:", apiURL)
	}

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("Error: DB_DSN environment variable not set.")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // Default port
	}
	// ... use these variables to initialize your adapters and server
	// For example, when starting Gin:
	// router.Run(":" + serverPort)
}
