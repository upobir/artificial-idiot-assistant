package ai

type Role int

const (
	System Role = iota
	User
	Ai
)

type Message struct {
	Role    Role
	Content string
}

// Raw Conversation
type Conversation struct {
	Messages []Message
}

func NewConversation() *Conversation {
	return &Conversation{
		Messages: []Message{},
	}
}

func (conv *Conversation) AddMessage(role Role, content string) {
	conv.Messages = append(conv.Messages, Message{
		Role:    role,
		Content: content,
	})
}
