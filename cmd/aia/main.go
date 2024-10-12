package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/upobir/artificial-idiot-assistant/internal/console"
	"github.com/upobir/artificial-idiot-assistant/internal/utils"
)

func main() {
	utils.InitializeLog()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	console.Run()
}
