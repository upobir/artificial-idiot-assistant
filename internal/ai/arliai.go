package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Arliai struct {
	Url    string
	ApiKey string
	Client *http.Client
}

func InitializeArliaiConfig(apiKey string) *Arliai {
	return &Arliai{
		Url:    "https://api.arliai.com/v1/chat/completions",
		ApiKey: apiKey,
		Client: &http.Client{},
	}
}

func (arliai *Arliai) ChatComplete(conv *Conversation) (Message, error) {
	payload, err := json.Marshal(conv)
	if err != nil {
		return Message{}, err
	}

	req, err := http.NewRequest("POST", arliai.Url, bytes.NewBuffer(payload))
	if err != nil {
		return Message{}, err
	}

	req.Header.Set("Authorization", "Bearer "+arliai.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := arliai.Client.Do(req)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Message{}, fmt.Errorf("failed request, status: %d, response: %v", resp.StatusCode, respBody)
	}

	var result struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return Message{}, err
	}

	if len(result.Choices) != 1 {
		return Message{}, fmt.Errorf("length mismatch with response choices: %v", result)
	}

	return result.Choices[0].Message, nil
}
