package taskqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type processor struct {
	q *db.Queries
}

func newProcessor(q *db.Queries) processor {
	return processor{
		q: q,
	}
}

func (p *processor) doProcessTaskCreateSubmissionRecord(
	ctx context.Context,
	payload *TaskPayloadCreateSubmissionRecord,
) (*TaskResultCreateSubmissionRecord, error) {
	// TODO: upsert
	submissionID, err := p.q.CreateSubmission(ctx, db.CreateSubmissionParams{
		GameID:   int32(payload.GameID()),
		UserID:   int32(payload.UserID()),
		Code:     payload.Code(),
		CodeSize: int32(payload.CodeSize),
	})
	if err != nil {
		return nil, err
	}

	return &TaskResultCreateSubmissionRecord{
		TaskPayload:  payload,
		SubmissionID: int(submissionID),
	}, nil
}

func (p *processor) doProcessTaskCompileSwiftToWasm(
	_ context.Context,
	payload *TaskPayloadCompileSwiftToWasm,
) (*TaskResultCompileSwiftToWasm, error) {
	type swiftcRequestData struct {
		MaxDuration int    `json:"max_duration_ms"`
		Code        string `json:"code"`
	}
	type swiftcResponseData struct {
		Status string `json:"status"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
	reqData := swiftcRequestData{
		MaxDuration: 5000,
		Code:        payload.Code(),
	}
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	res, err := http.Post("http://worker:80/api/swiftc", "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("http.Post failed: %v", err)
	}
	resData := swiftcResponseData{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, fmt.Errorf("json.Decode failed: %v", err)
	}
	return &TaskResultCompileSwiftToWasm{
		TaskPayload: payload,
		Status:      resData.Status,
		Stdout:      resData.Stdout,
		Stderr:      resData.Stderr,
	}, nil
}

func (p *processor) doProcessTaskCompileWasmToNativeExecutable(
	_ context.Context,
	payload *TaskPayloadCompileWasmToNativeExecutable,
) (*TaskResultCompileWasmToNativeExecutable, error) {
	type wasmcRequestData struct {
		MaxDuration int    `json:"max_duration_ms"`
		Code        string `json:"code"`
	}
	type wasmcResponseData struct {
		Status string `json:"status"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
	reqData := wasmcRequestData{
		MaxDuration: 5000,
		Code:        payload.Code(),
	}
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	res, err := http.Post("http://worker:80/api/wasmc", "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("http.Post failed: %v", err)
	}
	resData := wasmcResponseData{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, fmt.Errorf("json.Decode failed: %v", err)
	}
	return &TaskResultCompileWasmToNativeExecutable{
		TaskPayload: payload,
		Status:      resData.Status,
		Stdout:      resData.Stdout,
		Stderr:      resData.Stderr,
	}, nil
}

func (p *processor) doProcessTaskRunTestcase(
	_ context.Context,
	payload *TaskPayloadRunTestcase,
) (*TaskResultRunTestcase, error) {
	type testrunRequestData struct {
		MaxDuration int    `json:"max_duration_ms"`
		Code        string `json:"code"`
		Stdin       string `json:"stdin"`
	}
	type testrunResponseData struct {
		Status string `json:"status"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
	reqData := testrunRequestData{
		MaxDuration: 5000,
		Code:        payload.Code(),
		Stdin:       payload.Stdin,
	}
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	res, err := http.Post("http://worker:80/api/testrun", "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("http.Post failed: %v", err)
	}
	resData := testrunResponseData{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, fmt.Errorf("json.Decode failed: %v", err)
	}
	return &TaskResultRunTestcase{
		TaskPayload: payload,
		Status:      resData.Status,
		Stdout:      resData.Stdout,
		Stderr:      resData.Stderr,
	}, nil
}
