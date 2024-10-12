package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/upobir/artificial-idiot-assistant/internal/console"
	"github.com/upobir/artificial-idiot-assistant/internal/db"
	"github.com/upobir/artificial-idiot-assistant/internal/utils"
)

func main() {
	// log setup
	utils.InitializeLog()

	// env setup
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	// mongo setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := db.MongoConfig{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
	}
	client, err := db.MongoConnect(ctx, config)
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("DB disconnect failed: %v", err)
		}
	}()

	// console run
	console.Run()
}
