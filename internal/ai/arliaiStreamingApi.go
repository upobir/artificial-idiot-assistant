package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
		payload, err := json.Marshal(arliaiConversation{
			Model:             arliai.model,
			Messages:          mapMessages(conv.Messages),
			RepetitionPenalty: 1.1,
			Temperature:       0.5,
			TopP:              0.9,
			TopK:              40,
			MaxTokens:         300,
			Stream:            true,
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

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			ch <- ChatPart{Value: "", Err: fmt.Errorf("failed request, status: %d, response: %v", resp.StatusCode, string(body))}
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()

			if len(line) == 0 {
				continue
			}

			if !strings.HasPrefix(line, "data: ") {
				ch <- ChatPart{Value: "", Err: fmt.Errorf("unexpected line: %s", line)}
				return
			}

			line = strings.TrimPrefix(line, "data: ")

			if line == "[DONE]" {
				return
			}

			var result arliaiStreamingResponse
			if err := json.Unmarshal([]byte(line), &result); err != nil {
				ch <- ChatPart{Value: "", Err: err}
				return
			}

			if len(result.Choices) != 1 {
				ch <- ChatPart{Value: "", Err: fmt.Errorf("length mismatch with response choices: %v", result)}
				return
			}

			if result.Choices[0].Delta.Role != "assistant" && result.Choices[0].Delta.Role != "" {
				ch <- ChatPart{Value: "", Err: fmt.Errorf("unexpected role: %v", result)}
				return
			}

			ch <- ChatPart{Value: result.Choices[0].Delta.Content, Err: nil}
		}
	}()

	return ch
}
