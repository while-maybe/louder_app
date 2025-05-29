package main

import (
	"log"
	dbdriven "louder/internal/adapters/driven/mock_db"
	apidriving "louder/internal/adapters/driving/api_provider/stdlib"
	coreservice "louder/internal/core/service"
	"louder/pkg/config"
)

func main() {

	cfg := config.LoadConfig()

	log.Println("LOUDER starting")

	// serverAddr := os.Getenv("SERVER_ADDR")

	// instantiate driven adapters
	dataRepo := dbdriven.NewMockDBMessageRepository("Message With Time")
	randomRepo := dbdriven.NewMockDBRandRepository("Random number")

	// instantiate core app services
	messageService := coreservice.NewMessageService(dataRepo)
	randomNumberService := coreservice.NewRandNumberService(randomRepo)

	// instantiate driving adapters
	messageHandler := apidriving.NewMessageHandler(messageService)
	randomNumberHandler := apidriving.NewRandomNumberHandler(randomNumberService)

	// instantiate router
	router := apidriving.NewRouter(messageHandler, randomNumberHandler)

	// start the server
	log.Printf("Starting server on port %s\n", cfg.ServerPort)
	apidriving.StartServer(router, cfg.ServerPort)
	log.Println("Shutting server down")
}
