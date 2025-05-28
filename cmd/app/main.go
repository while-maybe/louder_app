package main

import (
	"log"
	dbdriven "louder/internal/adapters/driven/mock_db"
	apidriving "louder/internal/adapters/driving/api_provider/stdlib"
	coreservice "louder/internal/core/service"
	"os"
)

func main() {
	log.Print("LOUDER: ")

	serverAddr := os.Getenv("SERVER_ADDR")

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
	log.Printf("Starting server on port %s\n", serverAddr)
	apidriving.StartServer(router, serverAddr)
	log.Println("Shutting server down")
}
