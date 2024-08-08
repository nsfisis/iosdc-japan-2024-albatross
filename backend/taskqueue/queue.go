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
	codeHash MD5HexHash,
) error {
	task, err := newTaskCreateSubmissionRecord(
		gameID,
		userID,
		code,
		codeSize,
		codeHash,
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
	codeHash MD5HexHash,
	submissionID int,
) error {
	task, err := newTaskCompileSwiftToWasm(
		gameID,
		userID,
		code,
		codeHash,
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
	codeHash MD5HexHash,
	submissionID int,
) error {
	task, err := newTaskCompileWasmToNativeExecutable(
		gameID,
		userID,
		codeHash,
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
	codeHash MD5HexHash,
	submissionID int,
	testcaseID int,
	stdin string,
	stdout string,
) error {
	task, err := newTaskRunTestcase(
		gameID,
		userID,
		codeHash,
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
