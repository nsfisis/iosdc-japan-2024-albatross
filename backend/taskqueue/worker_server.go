package taskqueue

import (
	"github.com/hibiken/asynq"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type WorkerServer struct {
	server    *asynq.Server
	processor *processor
}

func NewWorkerServer(redisAddr string, queries *db.Queries) *WorkerServer {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{},
	)
	processor := newProcessor(queries)
	return &WorkerServer{
		server:    server,
		processor: processor,
	}
}

func (s *WorkerServer) Run() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(string(TaskTypeCreateSubmissionRecord), s.processor.processTaskCreateSubmissionRecord)
	mux.HandleFunc(string(TaskTypeCompileSwiftToWasm), s.processor.processTaskCompileSwiftToWasm)
	mux.HandleFunc(string(TaskTypeCompileWasmToNativeExecutable), s.processor.processTaskCompileWasmToNativeExecutable)
	mux.HandleFunc(string(TaskTypeRunTestcase), s.processor.processTaskRunTestcase)

	return s.server.Run(mux)
}

func (s *WorkerServer) Results() chan TaskResult {
	return s.processor.results
}
