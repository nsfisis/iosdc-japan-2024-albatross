package taskqueue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type TaskType string

// MD5 hash in hexadecimal format
type MD5HexHash string

const (
	TaskTypeCreateSubmissionRecord        TaskType = "create_submission_record"
	TaskTypeCompileSwiftToWasm            TaskType = "compile_swift_to_wasm"
	TaskTypeCompileWasmToNativeExecutable TaskType = "compile_wasm_to_native_executable"
	TaskTypeRunTestcase                   TaskType = "run_testcase"
)

type TaskPayloadBase struct {
	GameID   int
	UserID   int
	CodeHash MD5HexHash
}

type TaskPayloadCreateSubmissionRecord struct {
	TaskPayloadBase
	Code     string
	CodeSize int
}

func newTaskCreateSubmissionRecord(
	gameID int,
	userID int,
	code string,
	codeSize int,
	codeHash MD5HexHash,
) (*asynq.Task, error) {
	payload, err := json.Marshal(TaskPayloadCreateSubmissionRecord{
		TaskPayloadBase: TaskPayloadBase{
			GameID:   gameID,
			UserID:   userID,
			CodeHash: codeHash,
		},
		Code:     code,
		CodeSize: codeSize,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(TaskTypeCreateSubmissionRecord), payload), nil
}

func (t *TaskPayloadCreateSubmissionRecord) GameID() int          { return t.TaskPayloadBase.GameID }
func (t *TaskPayloadCreateSubmissionRecord) UserID() int          { return t.TaskPayloadBase.UserID }
func (t *TaskPayloadCreateSubmissionRecord) CodeHash() MD5HexHash { return t.TaskPayloadBase.CodeHash }

type TaskPayloadCompileSwiftToWasm struct {
	TaskPayloadBase
	SubmissionID int
	Code         string
}

func newTaskCompileSwiftToWasm(
	gameID int,
	userID int,
	code string,
	codeHash MD5HexHash,
	submissionID int,
) (*asynq.Task, error) {
	payload, err := json.Marshal(TaskPayloadCompileSwiftToWasm{
		TaskPayloadBase: TaskPayloadBase{
			GameID:   gameID,
			UserID:   userID,
			CodeHash: codeHash,
		},
		SubmissionID: submissionID,
		Code:         code,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(TaskTypeCompileSwiftToWasm), payload), nil
}

func (t *TaskPayloadCompileSwiftToWasm) GameID() int          { return t.TaskPayloadBase.GameID }
func (t *TaskPayloadCompileSwiftToWasm) UserID() int          { return t.TaskPayloadBase.UserID }
func (t *TaskPayloadCompileSwiftToWasm) CodeHash() MD5HexHash { return t.TaskPayloadBase.CodeHash }

type TaskPayloadCompileWasmToNativeExecutable struct {
	TaskPayloadBase
	SubmissionID int
}

func newTaskCompileWasmToNativeExecutable(
	gameID int,
	userID int,
	codeHash MD5HexHash,
	submissionID int,
) (*asynq.Task, error) {
	payload, err := json.Marshal(TaskPayloadCompileWasmToNativeExecutable{
		TaskPayloadBase: TaskPayloadBase{
			GameID:   gameID,
			UserID:   userID,
			CodeHash: codeHash,
		},
		SubmissionID: submissionID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(TaskTypeCompileWasmToNativeExecutable), payload), nil
}

func (t *TaskPayloadCompileWasmToNativeExecutable) GameID() int { return t.TaskPayloadBase.GameID }
func (t *TaskPayloadCompileWasmToNativeExecutable) UserID() int { return t.TaskPayloadBase.UserID }
func (t *TaskPayloadCompileWasmToNativeExecutable) CodeHash() MD5HexHash {
	return t.TaskPayloadBase.CodeHash
}

type TaskPayloadRunTestcase struct {
	TaskPayloadBase
	SubmissionID int
	TestcaseID   int
	Stdin        string
	Stdout       string
}

func newTaskRunTestcase(
	gameID int,
	userID int,
	codeHash MD5HexHash,
	submissionID int,
	testcaseID int,
	stdin string,
	stdout string,
) (*asynq.Task, error) {
	payload, err := json.Marshal(TaskPayloadRunTestcase{
		TaskPayloadBase: TaskPayloadBase{
			GameID:   gameID,
			UserID:   userID,
			CodeHash: codeHash,
		},
		SubmissionID: submissionID,
		TestcaseID:   testcaseID,
		Stdin:        stdin,
		Stdout:       stdout,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(TaskTypeRunTestcase), payload), nil
}

func (t *TaskPayloadRunTestcase) GameID() int          { return t.TaskPayloadBase.GameID }
func (t *TaskPayloadRunTestcase) UserID() int          { return t.TaskPayloadBase.UserID }
func (t *TaskPayloadRunTestcase) CodeHash() MD5HexHash { return t.TaskPayloadBase.CodeHash }

type TaskResult interface {
	Type() TaskType
	GameID() int
}

type TaskResultCreateSubmissionRecord struct {
	TaskPayload  *TaskPayloadCreateSubmissionRecord
	SubmissionID int
	Err          error
}

func (r *TaskResultCreateSubmissionRecord) Type() TaskType { return TaskTypeCreateSubmissionRecord }
func (r *TaskResultCreateSubmissionRecord) GameID() int    { return r.TaskPayload.GameID() }

type TaskResultCompileSwiftToWasm struct {
	TaskPayload *TaskPayloadCompileSwiftToWasm
	Status      string
	Stdout      string
	Stderr      string
	Err         error
}

func (r *TaskResultCompileSwiftToWasm) Type() TaskType { return TaskTypeCompileSwiftToWasm }
func (r *TaskResultCompileSwiftToWasm) GameID() int    { return r.TaskPayload.GameID() }

type TaskResultCompileWasmToNativeExecutable struct {
	TaskPayload *TaskPayloadCompileWasmToNativeExecutable
	Status      string
	Stdout      string
	Stderr      string
	Err         error
}

func (r *TaskResultCompileWasmToNativeExecutable) Type() TaskType {
	return TaskTypeCompileWasmToNativeExecutable
}
func (r *TaskResultCompileWasmToNativeExecutable) GameID() int { return r.TaskPayload.GameID() }

type TaskResultRunTestcase struct {
	TaskPayload *TaskPayloadRunTestcase
	Status      string
	Stdout      string
	Stderr      string
	Err         error
}

func (r *TaskResultRunTestcase) Type() TaskType { return TaskTypeRunTestcase }
func (r *TaskResultRunTestcase) GameID() int    { return r.TaskPayload.GameID() }
