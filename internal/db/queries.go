package db

import (
	"context"
	"errors"

	"github.com/upobir/artificial-idiot-assistant/internal/task"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const TASK_COLLECTION_NAME = "tasks"

func GetAllTasks(database *mongo.Database) ([]task.Task, error) {
	collection := database.Collection(TASK_COLLECTION_NAME)
	var tasks []task.Task
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func InsertTask(database *mongo.Database, task task.Task) (task.Task, error) {
	collection := database.Collection(TASK_COLLECTION_NAME)

	nextID, err := GetNextTaskId(database)
	if err != nil {
		return task, err
	}

	task.ID = nextID
	_, err = collection.InsertOne(context.TODO(), task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func UpdateTask(database *mongo.Database, task task.Task) (task.Task, error) {
	collection := database.Collection(TASK_COLLECTION_NAME)

	filter := bson.M{"id": task.ID}
	update := bson.M{
		"$set": bson.M{
			"name":  task.Name,
			"state": task.State,
			"tags":  task.Tags,
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return task, err
	}
	if result.MatchedCount == 0 {
		return task, errors.New("no task found with the given id")
	}
	return task, nil
}

func DeleteTask(database *mongo.Database, id int) error {
	collection := database.Collection(TASK_COLLECTION_NAME)

	filter := bson.M{"id": id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no task found with the given id")
	}

	return nil
}
