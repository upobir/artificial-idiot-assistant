package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNextTaskId(database *mongo.Database) (int, error) {
	collection := database.Collection("sequences")
	var result struct {
		Value int `bson:"value"`
	}
	filter := bson.M{"_id": "tasks"}
	update := bson.M{"$inc": bson.M{"value": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Value, nil
}
