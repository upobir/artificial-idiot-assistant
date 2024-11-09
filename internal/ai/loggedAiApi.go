package ai

import (
	"log"
	"strings"
	"time"

	"github.com/upobir/artificial-idiot-assistant/internal/conversationLog"
	"github.com/upobir/artificial-idiot-assistant/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoggedAiApi struct {
	aiApi    loggableAiApi
	database *mongo.Database
	logId    int
	metadata map[string]any
}

func InitializeLoggedAiApi(aiApi loggableAiApi, database *mongo.Database) *LoggedAiApi {

	logId, err := db.GetNextId(database, "conversation-logs")
	if err != nil {
		log.Printf("could not get id for conversation logs\n")
	}

	return &LoggedAiApi{
		aiApi:    aiApi,
		database: database,
		logId:    logId,
		metadata: aiApi.Metadata(),
	}
}

func logFromConversation(conv *Conversation, id int, metadata map[string]any) *conversationLog.ConversationLog {
	messages := make([]conversationLog.MessageLog, len(conv.Messages))
	for index, message := range conv.Messages {
		messages[index].Content = message.Content
		switch message.Role {
		case System:
			messages[index].Role = "prompt"
		case User:
			messages[index].Role = "user"
		case Ai:
			messages[index].Role = "ai"
		}
	}

	return &conversationLog.ConversationLog{
		ID:        id,
		Timestamp: time.Now(),
		Messages:  messages,
		Metadata:  metadata,
	}
}

func (ai *LoggedAiApi) ChatComplete(conv *Conversation) <-chan ChatPart {
	ch := make(chan ChatPart)

	convLog := logFromConversation(conv, ai.logId, ai.metadata)
	var content strings.Builder
	convLog.Messages = append(convLog.Messages, conversationLog.MessageLog{Role: "ai"})

	go func() {
		defer close(ch)
		for part := range ai.aiApi.ChatComplete(conv) {
			if part.Err != nil {
				convLog.Error = part.Err.Error()
			} else {
				content.WriteString(part.Value)
			}
			ch <- part
		}

		convLog.Messages[len(convLog.Messages)-1].Content = content.String()

		if ai.logId > 0 {
			err := db.UpsertConversationLog(ai.database, convLog)
			if err != nil {
				log.Printf("could not save conversation log\n")
			}
		}
	}()

	return ch
}
