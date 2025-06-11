package main

import (
	"context"
	"log"
	sqlitedbadapter "louder/internal/adapters/driven/db"
	bunadapter "louder/internal/adapters/driven/db/bun_adapter"
	randomgenerator "louder/internal/adapters/driven/random_generator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// sqlxadapter "louder/internal/adapters/driven/db/sqlx_adapter"
	dbdriven "louder/internal/adapters/driven/mock_db"
	apidriving "louder/internal/adapters/driving/api_provider/stdlib"
	coreservice "louder/internal/core/service"
	randomnumber "louder/internal/core/service/randomnumbers"
	"louder/pkg/config"
)

func main() {

	cfg := config.LoadConfig()
	// serverAddr := os.Getenv("SERVER_ADDR")
	log.Println("LOUDER starting")

	// instantiate driven adapters
	dataRepo := dbdriven.NewMockDBMessageRepository("Message With Time")
	randomGen := randomgenerator.NewStdLibGenerator()

	// instantiate driven adapter for sqlitedb
	db, err := sqlitedbadapter.Init("./louder.db")
	if err != nil {
		log.Fatalf("error cannot init DB: %s", err)
	}
	defer func() {
		log.Println("closing DB connection...")
		if err := db.Close(); err != nil {
			log.Printf("error closing DB: %v", err)
		}
	}() // this is deferred to ensure it happens on shutdown

	// Define the path to your migration files
	migrationsPath := "./migrations" // Best to move this to your config eventually

	// Run database migrations instead of creating new one
	err = sqlitedbadapter.RunMigrations(db, migrationsPath)
	if err != nil {
		log.Fatalf("error cannot run database migrations: %v", err)
	}
	// err = sqlitedbadapter.CreateSchema(db)
	// if err != nil {
	// 	log.Fatalf("error cannot initialise db schema: %s", err)
	// }

	// a little silly at the moment but Creating a single person is done through Bun, the whole list of people through SQLx
	singlePostRepo, err := bunadapter.NewBunPersonRepo(db)
	if err != nil {
		log.Fatalf("error cannot instantiate DB via Bun")
	}
	// the rest is done via SQLx
	// peopleRepo, err := sqlxadapter.NewSQLxPersonRepo(db)
	// if err != nil {
	// 	log.Fatalf("error cannot instantiate DB via SQLx")
	// }

	// instantiate core app services
	messageService := coreservice.NewMessageService(dataRepo)
	randomNumberService := randomnumber.NewRandNumberService(randomGen)

	// instantiate single Person get via Bun
	singlePostService := coreservice.NewPersonService(singlePostRepo)
	// instantiate Person core app service
	// peopleService := coreservice.NewPersonService(peopleRepo)

	// instantiate driving adapters
	messageHandler := apidriving.NewMessageHandler(messageService)
	randomNumberHandler := apidriving.NewRandomNumberHandler(randomNumberService)

	// for now with only the POST user Handler
	singlePostHandler := apidriving.NewPersonHandler(singlePostService)

	// and everything else will go here
	// var _ *coreservice.personServiceImpl = peopleService

	// instantiate router
	router := apidriving.NewRouter(messageHandler, randomNumberHandler, singlePostHandler)

	// wrap the router in a timeout handler - every incoming request will have a 5 sec deadline
	timeoutDuration := 5 * time.Second
	timedHandler := http.TimeoutHandler(router, timeoutDuration, "request timed out")

	// gracefully shutdown
	stdAPIServer := apidriving.NewStdAPIServer(":"+cfg.ServerPort, timedHandler)

	// channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// we start the server in a non-blocking way (go routine)
	go func() {
		if err := stdAPIServer.ListenAndServe(); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	// block until a signal is received
	<-stopChan

	// invoke graceful shutdown
	log.Println("shutdown signal received. starting graceful shutdown...")

	// this is a new context with a timeout for the graceful shutdown only
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelShutdown()

	// do the actual shutdown here
	if err := stdAPIServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdwon failed :(")
	}

	log.Printf("server shutdown gracefully")
}
