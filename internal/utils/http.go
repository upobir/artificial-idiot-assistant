package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func PostJson(url string, apiKey string, data any, client *http.Client, result any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed request, status: %d, response: %v", resp.StatusCode, respBody)
	}

	if err := json.Unmarshal(respBody, result); err != nil {

		return err
	}

	return nil
}

func PostJsonAndConsumeSse(url string, apiKey string, data any, client *http.Client, result any, callback func(any) error) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			return fmt.Errorf("unexpected line: %s", line)
		}

		line = strings.TrimPrefix(line, "data: ")

		if line == "[DONE]" {
			return nil
		}

		if err := json.Unmarshal([]byte(line), result); err != nil {
			return err
		}

		if err := callback(result); err != nil {
			return err
		}
	}

	return nil
}
