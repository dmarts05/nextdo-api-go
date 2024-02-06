package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dmarts05/nextdo-api-go/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Setup
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, using default values")
	}

	s, err := server.New()
	if err != nil {
		log.Fatal(err)
	}

	// Interrupt signal to gracefully shutdown the server
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start server
	err = s.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
