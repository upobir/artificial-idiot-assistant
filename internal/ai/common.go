package ai

type apiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func mapMessages(messages []Message) []apiMessage {
	result := make([]apiMessage, len(messages))
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
	Model             string       `json:"model"`
	Messages          []apiMessage `json:"messages"`
	RepetitionPenalty float32      `json:"repetition_penalty"`
	Temperature       float32      `json:"temperature"`
	TopP              float32      `json:"top_p"`
	TopK              int          `json:"top_k"`
	MaxTokens         int          `json:"max_tokens"`
	Stream            bool         `json:"stream"`
}

type groqConversation struct {
	Model           string       `json:"model"`
	Messages        []apiMessage `json:"messages"`
	PresencePenalty float32      `json:"presence_penalty"`
	Temperature     float32      `json:"temperature"`
	TopP            float32      `json:"top_p"`
	MaxTokens       int          `json:"max_tokens"`
	Stream          bool         `json:"stream"`
}

type apiResponse struct {
	Choices []struct {
		Message apiMessage `json:"message"`
	} `json:"choices"`
}

type apiStreamingResponse struct {
	Choices []struct {
		Delta apiMessage `json:"delta"`
	} `json:"choices"`
}
