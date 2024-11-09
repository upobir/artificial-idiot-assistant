package db

import (
	"context"

	"github.com/upobir/artificial-idiot-assistant/internal/conversationLog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CONVERSATION_LOG_COLLECTION_NAME = "conversation-logs"

func UpsertConversationLog(database *mongo.Database, conv *conversationLog.ConversationLog) error {
	collection := database.Collection(CONVERSATION_LOG_COLLECTION_NAME)

	filter := bson.M{"id": conv.ID}
	update := bson.M{
		"$set": *conv,
	}

	updateOptions := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, updateOptions)
	return err
}
