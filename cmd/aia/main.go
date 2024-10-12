package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
	"github.com/upobir/artificial-idiot-assistant/internal/console"
	"github.com/upobir/artificial-idiot-assistant/internal/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	utils.InitializeLog()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:rootpassword@localhost:27017"))

	if err != nil {
		log.Fatalf("Client connect failed: %v", err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Client disconnect failed: %v", err)
		}
	}()
	{
		var result bson.M
		if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
			panic(err)
		}

		log.Printf("DB connected")
	}

	{
		db := client.Database("aia-dev")
		collection := db.Collection("test")

		testDocument := map[string]any{
			"name":      "test",
			"value":     123,
			"timestamp": time.Now(),
		}

		result, err := collection.InsertOne(ctx, testDocument)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted document ID:", result.InsertedID)
	}

	console.Run()
}
