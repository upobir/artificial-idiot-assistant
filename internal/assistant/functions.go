package assistant

import (
	"encoding/json"
	"strconv"

	"github.com/upobir/artificial-idiot-assistant/internal/db"
	"github.com/upobir/artificial-idiot-assistant/internal/task"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler func(database *mongo.Database, arguments string) (string, error)

var AssistantHandlers = map[string]Handler{
	"getTasks":   getTasks,
	"insertTask": insertTask,
	"updateTask": updateTask,
}

func getTasks(database *mongo.Database, arguments string) (string, error) {
	tasks, err := db.GetAllTasks(database)
	if err != nil {
		return "", err
	}
	result, err := json.Marshal(tasks)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func insertTask(database *mongo.Database, arguments string) (string, error) {
	var task task.Task
	if err := json.Unmarshal([]byte(arguments), &task); err != nil {
		return "", err
	}
	task, err := db.InsertTask(database, task)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(task.ID), nil
}

func updateTask(database *mongo.Database, arguments string) (string, error) {
	var task task.Task
	if err := json.Unmarshal([]byte(arguments), &task); err != nil {
		return "", err
	}
	task, err := db.UpdateTask(database, task)
	if err != nil {
		return "", err
	}
	return "true", err
}
