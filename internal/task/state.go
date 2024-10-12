package task

type TaskState string

const (
	SCRATCHPAD TaskState = "SCRATCHPAD"
	PENDING    TaskState = "PENDING"
	BLOCKED    TaskState = "BLOCKED"
	COMPLETED  TaskState = "COMPLETED"
	DISCARDED  TaskState = "DISCARDED"
)

func (s TaskState) String() string {
	return string(s)
}
