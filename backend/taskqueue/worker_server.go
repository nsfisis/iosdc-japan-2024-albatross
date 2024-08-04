package taskqueue

import (
	"github.com/hibiken/asynq"
)

type WorkerServer struct {
	server *asynq.Server
}

func NewWorkerServer(redisAddr string) *WorkerServer {
	return &WorkerServer{
		server: asynq.NewServer(
			asynq.RedisClientOpt{
				Addr: redisAddr,
			},
			asynq.Config{},
		),
	}
}

func (s *WorkerServer) Run() error {
	mux := asynq.NewServeMux()
	mux.Handle(TaskTypeExec, NewExecProcessor())

	return s.server.Run(mux)
}
