package ai

import (
	"fmt"
	"net/http"

	"github.com/upobir/artificial-idiot-assistant/internal/utils"
)

type ArliaiApi struct {
	url    string
	apiKey string
	client *http.Client
	model  string
}

func InitializeArliaiApi(apiKey string, model string) *ArliaiApi {
	return &ArliaiApi{
		url:    "https://api.arliai.com/v1/chat/completions",
		apiKey: apiKey,
		client: &http.Client{},
		model:  model,
	}
}

func (arliai *ArliaiApi) ChatComplete(conv *Conversation) <-chan ChatPart {
	ch := make(chan ChatPart)

	go func() {
		defer close(ch)
		paylod := arliaiConversation{
			Model:             arliai.model,
			Messages:          mapMessages(conv.Messages),
			RepetitionPenalty: 1.1,
			Temperature:       0.5,
			TopP:              0.9,
			TopK:              40,
			MaxTokens:         300,
			Stream:            false,
		}
		var result apiResponse
		err := utils.PostJson(arliai.url, arliai.apiKey, paylod, arliai.client, &result)
		if err != nil {
			ch <- ChatPart{Value: "", Err: err}
			return
		}

		if len(result.Choices) != 1 {
			ch <- ChatPart{Value: "", Err: fmt.Errorf("length mismatch with response choices: %v", result)}
			return
		}

		if result.Choices[0].Message.Role != "assistant" {
			ch <- ChatPart{Value: "", Err: fmt.Errorf("unexpected role: %v", result)}
			return
		}

		ch <- ChatPart{Value: result.Choices[0].Message.Content, Err: nil}
	}()

	return ch
}
