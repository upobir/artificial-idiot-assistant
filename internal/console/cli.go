package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/upobir/artificial-idiot-assistant/internal/ai"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run(database *mongo.Database, arliai *ai.Arliai) error {
	reader := bufio.NewReader(os.Stdin)

	conv := ai.NewConversation("Llama-3.1-8B-Storm", false)
	conv.AddLocalMessage("system", "You are a helpful assistant.")

	for {
		fmt.Print("you > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" {
			fmt.Println("Exiting...")
			break
		}

		conv.AddLocalMessage("user", input)
		msg, err := conv.FetchAssistantMessage(arliai, true)
		if err != nil {
			log.Fatalf("chat completion: %v\n", err)
		}
		if msg.Role != "assistant" {
			log.Fatalf("unexpected role: %v", msg.Role)
		}

		fmt.Printf("aia > %s\n", msg.Content)
		fmt.Println()
	}

	return nil
}
