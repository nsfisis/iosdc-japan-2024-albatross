package taskqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
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
		Code:     payload.Code,
		CodeSize: int32(payload.CodeSize),
		CodeHash: string(payload.CodeHash()),
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
		CodeHash    string `json:"code_hash"`
	}
	type swiftcResponseData struct {
		Status string `json:"status"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
	reqData := swiftcRequestData{
		MaxDuration: 5000,
		Code:        payload.Code,
		CodeHash:    string(payload.CodeHash()),
	}
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	req, err := http.NewRequest("POST", "http://worker:80/api/swiftc", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	jwt, err := auth.NewAnonymousJWT()
	if err != nil {
		return nil, fmt.Errorf("auth.NewAnonymousJWT failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do failed: %v", err)
	}
	defer res.Body.Close()

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
		CodeHash    string `json:"code_hash"`
	}
	type wasmcResponseData struct {
		Status string `json:"status"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
	reqData := wasmcRequestData{
		MaxDuration: 5000,
		CodeHash:    string(payload.CodeHash()),
	}
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	req, err := http.NewRequest("POST", "http://worker:80/api/wasmc", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	jwt, err := auth.NewAnonymousJWT()
	if err != nil {
		return nil, fmt.Errorf("auth.NewAnonymousJWT failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do failed: %v", err)
	}
	defer res.Body.Close()

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
		CodeHash    string `json:"code_hash"`
		Stdin       string `json:"stdin"`
	}
	type testrunResponseData struct {
		Status string `json:"status"`
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
	reqData := testrunRequestData{
		MaxDuration: 5000,
		CodeHash:    string(payload.CodeHash()),
		Stdin:       payload.Stdin,
	}
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %v", err)
	}
	req, err := http.NewRequest("POST", "http://worker:80/api/testrun", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	jwt, err := auth.NewAnonymousJWT()
	if err != nil {
		return nil, fmt.Errorf("auth.NewAnonymousJWT failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do failed: %v", err)
	}
	defer res.Body.Close()

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
