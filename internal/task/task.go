package task

type Task struct {
	ID    int       `bson:"id"`
	Name  string    `bson:"name"`
	State TaskState `bson:"state"`
	Tags  []string  `bson:"tag"`
}
