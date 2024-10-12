package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (config *MongoConfig) getUrl() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Username, config.Password, config.Host, config.Port)
}

func MongoConnect(ctx context.Context, config MongoConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.getUrl()))

	if err != nil {
		return nil, err
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}

	return client, nil
}
