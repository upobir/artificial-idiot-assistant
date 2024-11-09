package ai

import (
	"fmt"
	"net/http"

	"github.com/upobir/artificial-idiot-assistant/internal/utils"
)

type ArliaiStreamingApi struct {
	url    string
	apiKey string
	client *http.Client
	model  string
}

func InitializeArliaiStreamingApi(apiKey string, model string) *ArliaiStreamingApi {
	return &ArliaiStreamingApi{
		url:    "https://api.arliai.com/v1/chat/completions",
		apiKey: apiKey,
		client: &http.Client{},
		model:  model,
	}
}

type arliaiStreamingResponse struct {
	Choices []struct {
		Delta arliaiMessage `json:"delta"`
	} `json:"choices"`
}

func (arliai *ArliaiStreamingApi) ChatComplete(conv *Conversation) <-chan ChatPart {
	ch := make(chan ChatPart)

	go func() {
		defer close(ch)
		payload := arliaiConversation{
			Model:             arliai.model,
			Messages:          mapMessages(conv.Messages),
			RepetitionPenalty: 1.1,
			Temperature:       0.5,
			TopP:              0.9,
			TopK:              40,
			MaxTokens:         300,
			Stream:            true,
		}

		err := utils.PostJsonAndConsumeSse(arliai.url, arliai.apiKey, payload, arliai.client, &arliaiStreamingResponse{}, func(chunk any) error {
			response := chunk.(*arliaiStreamingResponse)

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
