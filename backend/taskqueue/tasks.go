package taskqueue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskTypeExec = "exec"
)

type TaskExecPlayload struct {
	GameID int
	UserID int
	Code   string
}

func NewExecTask(gameID, userID int, code string) (*asynq.Task, error) {
	payload, err := json.Marshal(TaskExecPlayload{
		GameID: gameID,
		UserID: userID,
		Code:   code,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTypeExec, payload), nil
}

type TaskExecResult struct {
	Task   *TaskExecPlayload
	Result string
	Stdout string
	Stderr string
}
