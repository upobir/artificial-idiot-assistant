package ai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	RepetitionPenalty float32   `json:"repetition_penalty"`
	Temperature       float32   `json:"temperature"`
	TopP              float32   `json:"top_p"`
	TopK              int       `json:"top_k"`
	MaxTokens         int       `json:"max_tokens"`
	Stream            bool      `json:"stream"`
}

func NewConversation(model string, stream bool) *Conversation {
	return &Conversation{
		Model:             model,
		Messages:          []Message{},
		RepetitionPenalty: 1.1,
		Temperature:       0.5,
		TopP:              0.9,
		TopK:              40,
		MaxTokens:         200,
		Stream:            stream,
	}
}

func (conv *Conversation) AddLocalMessage(role string, message string) {
	conv.Messages = append(conv.Messages, Message{
		Role:    role,
		Content: message,
	})
}

func (conv *Conversation) FetchAssistantMessage(arliai *Arliai, update bool) (Message, error) {
	message, err := arliai.ChatComplete(conv)
	if err != nil {
		return Message{}, err
	}

	if update {
		conv.AddLocalMessage(message.Role, message.Content)
	}

	return message, nil
}
