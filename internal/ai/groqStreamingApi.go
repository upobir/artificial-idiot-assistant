package ai

import (
	"fmt"
	"net/http"

	"github.com/upobir/artificial-idiot-assistant/internal/utils"
)

type GroqStreamingApi struct {
	url    string
	apiKey string
	client *http.Client
	model  string
}

func InitializeGroqStreamingApi(apiKey string, model string) *GroqStreamingApi {
	return &GroqStreamingApi{
		url:    "https://api.groq.com/openai/v1/chat/completions",
		apiKey: apiKey,
		client: &http.Client{},
		model:  model,
	}
}

func (groqai *GroqStreamingApi) ChatComplete(conv *Conversation) <-chan ChatPart {
	ch := make(chan ChatPart)

	go func() {
		defer close(ch)
		payload := groqConversation{
			Model:           groqai.model,
			Messages:        mapMessages(conv.Messages),
			PresencePenalty: 1.1,
			Temperature:     0.5,
			TopP:            0.9,
			MaxTokens:       300,
			Stream:          true,
		}

		err := utils.PostJsonAndConsumeSse(groqai.url, groqai.apiKey, payload, groqai.client, &apiStreamingResponse{}, func(chunk any) error {
			response := chunk.(*apiStreamingResponse)

			if len(response.Choices) != 1 {
				return fmt.Errorf("length mismatch with response choices: %v", response)
			}

			if response.Choices[0].Delta.Role != "assistant" && response.Choices[0].Delta.Role != "" {
				return fmt.Errorf("unexpected role: %v", response)
			}

			ch <- ChatPart{Value: response.Choices[0].Delta.Content, Err: nil}
			response.Choices = nil

			return nil
		})
		if err != nil {
			ch <- ChatPart{Value: "", Err: err}
			return
		}
	}()

	return ch
}
