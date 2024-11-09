package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/upobir/artificial-idiot-assistant/internal/ai"
	"github.com/upobir/artificial-idiot-assistant/internal/assistant"
	"github.com/upobir/artificial-idiot-assistant/internal/console"
	"github.com/upobir/artificial-idiot-assistant/internal/db"
	"github.com/upobir/artificial-idiot-assistant/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	initializeLogs()

	initializeEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, database := initializeMongo(ctx)
	defer client.Disconnect(ctx)

	aiApi := initializeAiApi()

	assistant := initializeAssistant(aiApi, database)

	if err := console.Run(assistant); err != nil {
		log.Fatalf("Run error: %v", err)
	}
}

func initializeAssistant(aiApi ai.AiApi, database *mongo.Database) *assistant.Assistant {
	assistant, err := assistant.NewAssistant(aiApi, database)
	if err != nil {
		log.Fatalf("Error creating assistant: %v\n", err)
	}
	return assistant
}

func initializeAiApi() ai.AiApi {
	return ai.InitializeArliaiApi(os.Getenv("ARLIAI_API_KEY"), "Mistral-Nemo-12B-SauerkrautLM")
	// return ai.InitializeFakeApi(true, 100*time.Millisecond)
}

func initializeEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}
}

func initializeLogs() {
	utils.InitializeLog()
}

func initializeMongo(ctx context.Context) (*mongo.Client, *mongo.Database) {
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

	return client, database
}
