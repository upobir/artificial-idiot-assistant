package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/upobir/artificial-idiot-assistant/internal/ai"
	"github.com/upobir/artificial-idiot-assistant/internal/assistant"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run(database *mongo.Database, arliai *ai.Arliai) error {
	reader := bufio.NewReader(os.Stdin)

	assistant, err := assistant.NewAssistant("Llama-3.1-8B-Fireplace2", false)
	if err != nil {
		return err
	}

	for {
		fmt.Print("you > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		startTime := time.Now()

		if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			fmt.Println("Exiting...")
			break
		}

		assistant.AddUserMessage(input)
		msg, err := assistant.FetchAssistantMessage(arliai, database)
		if err != nil {
			log.Fatalf("chat completion: %v\n", err)
		}

		output := strings.ReplaceAll(msg, "\n", "\n    ")

		endTime := time.Now()
		fmt.Printf("aia > %s\n", output)
		fmt.Printf("(%v seconds)\n", endTime.Sub(startTime).Seconds())
		fmt.Println()
	}

	return nil
}
