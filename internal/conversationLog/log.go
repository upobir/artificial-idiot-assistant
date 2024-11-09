package conversationLog

import "time"

type MessageLog struct {
	Role    string `bson:"role"`
	Content string `bson:"content"`
}

type ConversationLog struct {
	ID        int            `bson:"id"`
	Timestamp time.Time      `bson:"timestamp"`
	Metadata  map[string]any `bson:"metadata"`
	Messages  []MessageLog   `bson:"messages"`
	Error     string         `bson:"error"`
}
