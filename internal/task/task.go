package task

type Task struct {
	ID    int      `bson:"id"`
	Name  string   `bson:"name"`
	State string   `bson:"state"`
	Tags  []string `bson:"tag"`
}
