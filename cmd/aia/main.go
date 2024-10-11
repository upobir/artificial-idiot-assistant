package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("[%s] %s\n", time.Now().Format("2006/01/02 15:04:05"), string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	reader := bufio.NewReader(os.Stdin)

	log.Printf("Secret env value is %s\n", os.Getenv("TEST"))

	for {
		fmt.Print("you > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			fmt.Println("Exiting...")
			break
		}

		fmt.Printf("aia > %s\n", input)
		fmt.Println()
	}
}
