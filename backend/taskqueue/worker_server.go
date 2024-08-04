package taskqueue

import (
	"github.com/hibiken/asynq"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type WorkerServer struct {
	server  *asynq.Server
	queries *db.Queries
}

func NewWorkerServer(redisAddr string, queries *db.Queries) *WorkerServer {
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
	mux.Handle(TaskTypeExec, NewExecProcessor(s.queries))

	return s.server.Run(mux)
}
