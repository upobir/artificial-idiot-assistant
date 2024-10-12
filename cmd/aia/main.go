package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/upobir/artificial-idiot-assistant/internal/ai"
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

	dbConfig := db.MongoConfig{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
	}
	client, err := db.MongoConnect(ctx, dbConfig)
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}

	database := client.Database(os.Getenv("MONGO_DATABASE"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("DB disconnect failed: %v", err)
		}
	}()

	// arliai setup
	arliai := ai.InitializeArliaiConfig(os.Getenv("ARLIAI_API_KEY"))

	// console run
	if err = console.Run(database, arliai); err != nil {
		log.Fatalf("Run error: %v", err)
	}
}
