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

func (q *Queue) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return q.client.Enqueue(task, opts...)
}
