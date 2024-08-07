package taskqueue

import (
	"github.com/hibiken/asynq"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type WorkerServer struct {
	server  *asynq.Server
	queries *db.Queries
	results chan TaskExecResult
}

func NewWorkerServer(redisAddr string, queries *db.Queries) *WorkerServer {
	return &WorkerServer{
		server: asynq.NewServer(
			asynq.RedisClientOpt{
				Addr: redisAddr,
			},
			asynq.Config{},
		),
		queries: queries,
		results: make(chan TaskExecResult),
	}
}

func (s *WorkerServer) Run() error {
	mux := asynq.NewServeMux()
	mux.Handle(TaskTypeExec, NewExecProcessor(s.queries, s.results))

	return s.server.Run(mux)
}

func (s *WorkerServer) Results() chan TaskExecResult {
	return s.results
}
