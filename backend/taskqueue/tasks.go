package taskqueue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskTypeExec = "exec"
)

type Testcase struct {
	Stdin  string
	Stdout string
}

type TaskExecPlayload struct {
	Code      string
	Testcases []*Testcase
}

func NewExecTask(code string, testcases []*Testcase) (*asynq.Task, error) {
	payload, err := json.Marshal(TaskExecPlayload{
		Code:      code,
		Testcases: testcases,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTypeExec, payload), nil
}
