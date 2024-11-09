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

type Assistant struct {
	aiApi        ai.AiApi
	database     *mongo.Database
	conversation *ai.Conversation
}

func NewAssistant(aiApi ai.AiApi, database *mongo.Database) (*Assistant, error) {

	prompt, err := loadPrompt()
	if err != nil {
		return nil, err
	}

	assistant := &Assistant{
		aiApi:        aiApi,
		database:     database,
		conversation: ai.NewConversation(),
	}

	assistant.conversation.AddMessage(ai.System, prompt)

	return assistant, nil
}

func loadPrompt() (string, error) {
	promptPath := filepath.Join(os.Getenv("SYSTEM_PROMPT_PATH"), "system_prompt.txt")
	file, err := os.Open(promptPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (assistant *Assistant) AddUserMessage(message string) {
	assistant.conversation.AddMessage(ai.User, "FROM-USER: "+message)
}

type ResponseKind int

const (
	Action ResponseKind = iota
	MessagePart
	Unknown
)

type ResponsePart struct {
	Kind    ResponseKind
	Content string
	Err     error
}

func (assistant *Assistant) FetchAssistantMessage() <-chan ResponsePart {
	ch := make(chan ResponsePart)
	go func() {
		defer close(ch)
		for {
			var builder strings.Builder
			kind := Unknown
			for part := range assistant.aiApi.ChatComplete(assistant.conversation) {
				if part.Err != nil {
					ch <- ResponsePart{Kind: Unknown, Content: "", Err: part.Err}
					return
				}

				value := part.Value
				if builder.Len() == 0 {
					value = strings.TrimLeft(value, " ")
				}

				builder.WriteString(value)

				if kind == Unknown {
					message := builder.String()
					if strings.HasPrefix(message, "TO-USER: ") {
						kind = MessagePart
						content := strings.TrimPrefix(message, "TO-USER: ")
						ch <- ResponsePart{Kind: kind, Content: content, Err: nil}
					} else if strings.HasPrefix(message, "TO-SERVICE: ") {
						kind = Action
						ch <- ResponsePart{Kind: kind, Content: "Working...", Err: nil}
					}
				} else if kind == MessagePart {
					ch <- ResponsePart{Kind: kind, Content: value, Err: nil}
				} else if kind == Action {
				}
			}

			if kind == Unknown {
				ch <- ResponsePart{Kind: Unknown, Content: "", Err: fmt.Errorf("failed to parse: '%v'", builder.String())}
				return
			} else if kind == Action {
				response, err := assistant.handleFunctionCall(strings.TrimSpace(strings.TrimPrefix(builder.String(), "TO-SERVICE: ")))
				if err != nil {
					ch <- ResponsePart{Kind: Unknown, Content: "", Err: err}
					return
				}
				assistant.conversation.AddMessage(ai.User, "FROM-SERVICE: "+response)
			} else {
				return
			}
		}
	}()
	return ch
}

func (assistant *Assistant) handleFunctionCall(functionCall string) (string, error) {
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

	response, err := handler(assistant.database, arguments)
	if err != nil {
		return "", err
	}

	return response, nil
}
