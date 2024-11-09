package ai

type ChatPart struct {
	Value string
	Err   error
}

type AiApi interface {
	ChatComplete(conv *Conversation) <-chan ChatPart
}
