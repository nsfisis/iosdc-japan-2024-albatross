package taskqueue

import (
	"github.com/hibiken/asynq"
)

type Queue struct {
	client *asynq.Client
}

func NewQueue(redisAddr string) *Queue {
	return &Queue{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr: redisAddr,
		}),
	}
}

func (q *Queue) Close() {
	q.client.Close()
}

func (q *Queue) EnqueueTaskCreateSubmissionRecord(
	gameID int,
	userID int,
	code string,
	codeSize int,
) error {
	task, err := newTaskCreateSubmissionRecord(
		gameID,
		userID,
		code,
		codeSize,
	)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *Queue) EnqueueTaskCompileSwiftToWasm(
	gameID int,
	userID int,
	code string,
	submissionID int,
) error {
	task, err := newTaskCompileSwiftToWasm(
		gameID,
		userID,
		code,
		submissionID,
	)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *Queue) EnqueueTaskCompileWasmToNativeExecutable(
	gameID int,
	userID int,
	code string,
	submissionID int,
) error {
	task, err := newTaskCompileWasmToNativeExecutable(
		gameID,
		userID,
		code,
		submissionID,
	)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}

func (q *Queue) EnqueueTaskRunTestcase(
	gameID int,
	userID int,
	code string,
	submissionID int,
	testcaseID int,
	stdin string,
	stdout string,
) error {
	task, err := newTaskRunTestcase(
		gameID,
		userID,
		code,
		submissionID,
		testcaseID,
		stdin,
		stdout,
	)
	if err != nil {
		return err
	}
	_, err = q.client.Enqueue(task)
	return err
}
