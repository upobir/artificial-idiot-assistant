package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type arliaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func mapMessages(messages []Message) []arliaiMessage {
	result := make([]arliaiMessage, len(messages))
	for index, message := range messages {
		result[index].Content = message.Content
		switch message.Role {
		case System:
			result[index].Role = "system"
		case User:
			result[index].Role = "user"
		case Ai:
			result[index].Role = "assistant"
		}
	}
	return result
}

type arliaiConversation struct {
	Model             string          `json:"model"`
	Messages          []arliaiMessage `json:"messages"`
	RepetitionPenalty float32         `json:"repetition_penalty"`
	Temperature       float32         `json:"temperature"`
	TopP              float32         `json:"top_p"`
	TopK              int             `json:"top_k"`
	MaxTokens         int             `json:"max_tokens"`
	Stream            bool            `json:"stream"`
}

type arliaiResponse struct {
	Choices []struct {
		Message arliaiMessage `json:"message"`
	} `json:"choices"`
}

func (arliai *ArliaiApi) ChatComplete(conv *Conversation) <-chan ChatPart {
	ch := make(chan ChatPart)

	go func() {
		defer close(ch)
		payload, err := json.Marshal(arliaiConversation{
			Model:             arliai.model,
			Messages:          mapMessages(conv.Messages),
			RepetitionPenalty: 1.1,
			Temperature:       0.5,
			TopP:              0.9,
			TopK:              40,
			MaxTokens:         300,
			Stream:            false,
		})
		if err != nil {
			ch <- ChatPart{Value: "", Err: err}
			return
		}

		req, err := http.NewRequest("POST", arliai.url, bytes.NewBuffer(payload))
		if err != nil {
			ch <- ChatPart{Value: "", Err: err}
			return
		}

		req.Header.Set("Authorization", "Bearer "+arliai.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := arliai.client.Do(req)
		if err != nil {
			ch <- ChatPart{Value: "", Err: err}
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			ch <- ChatPart{Value: "", Err: err}
			return
		}

		if resp.StatusCode != http.StatusOK {
			ch <- ChatPart{Value: "", Err: fmt.Errorf("failed request, status: %d, response: %v", resp.StatusCode, respBody)}
			return
		}

		var result arliaiResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
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
