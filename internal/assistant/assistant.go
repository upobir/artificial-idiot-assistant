package assistant

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/upobir/artificial-idiot-assistant/internal/ai"
	"go.mongodb.org/mongo-driver/mongo"
)

type assistant struct {
	Conversation *ai.Conversation
}

func NewAssistant(model string, stream bool) (*assistant, error) {
	promptPath := filepath.Join(os.Getenv("SYSTEM_PROMPT_PATH"), "system_prompt.txt")
	file, err := os.Open(promptPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	assistant := &assistant{
		Conversation: ai.NewConversation("Mistral-Nemo-12B-Instruct-2407", false),
	}
	assistant.Conversation.AddLocalMessage("system", string(content))

	return assistant, nil
}

func (assistant *assistant) AddUserMessage(message string) {
	assistant.Conversation.AddLocalMessage("user", "FROM-USER: "+message)
}

func (assistant *assistant) FetchAssistantMessage(arliai *ai.Arliai, database *mongo.Database) (string, error) {
	for {
		msg, err := assistant.Conversation.FetchAssistantMessage(arliai, true)
		if err != nil {
			return "", err
		}
		if msg.Role != "assistant" {
			return "", fmt.Errorf("unknown role: %s", msg.Role)
		}

		if strings.HasPrefix(msg.Content, "TO-USER:") {
			return strings.TrimSpace(strings.TrimPrefix(msg.Content, "TO-USER:")), nil
		} else if strings.HasPrefix(msg.Content, "TO-SERVICE:") {
			response, err := assistant.handleFunctionCall(database, msg.Content)
			if err != nil {
				return "", err
			}
			assistant.Conversation.AddLocalMessage("user", "FROM-SERICE: "+response)
		} else {
			return "", fmt.Errorf("bad format response: %v", msg.Content)
		}
	}
}

func (assistant *assistant) handleFunctionCall(database *mongo.Database, message string) (string, error) {
	// log.Println(message)
	functionCall := strings.TrimSpace(strings.TrimPrefix(message, "TO-SERVICE:"))
	spaceInd := strings.Index(functionCall, " ")
	arguments := ""
	if spaceInd != -1 {
		arguments = functionCall[spaceInd:]
		functionCall = functionCall[:spaceInd]
	}
	functionCall = strings.TrimSpace(functionCall)
	arguments = strings.TrimSpace(arguments)

	handler, contains := AssistantHandlers[functionCall]
	if !contains {
		return "", fmt.Errorf("unknown function: %s", functionCall)
	}

	response, err := handler(database, arguments)
	if err != nil {
		return "", err
	}

	return response, nil
}
