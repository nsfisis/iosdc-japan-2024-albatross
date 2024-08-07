package taskqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type processor struct {
	q       *db.Queries
	results chan TaskResult
}

func newProcessor(q *db.Queries) *processor {
	return &processor{
		q:       q,
		results: make(chan TaskResult),
	}
}

func (p *processor) processTaskCreateSubmissionRecord(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadCreateSubmissionRecord
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		// TODO: send result
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// TODO: upsert
	// Create submission record.
	submissionID, err := p.q.CreateSubmission(ctx, db.CreateSubmissionParams{
		GameID:   int32(payload.GameID()),
		UserID:   int32(payload.UserID()),
		Code:     payload.Code(),
		CodeSize: int32(payload.CodeSize),
	})
	if err != nil {
		// TODO: send result
		return fmt.Errorf("CreateSubmission failed: %v", err)
	}

	p.results <- &TaskResultCreateSubmissionRecord{
		TaskPayload:  &payload,
		SubmissionID: int(submissionID),
	}
	return nil
}

func (p *processor) processTaskCompileSwiftToWasm(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadCompileSwiftToWasm
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		// TODO: send result
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

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
	reqJson, err := json.Marshal(reqData)
	if err != nil {
		// TODO: send result
		return fmt.Errorf("json.Marshal failed: %v", err)
	}
	res, err := http.Post("http://worker:80/api/swiftc", "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		// TODO: send result
		return fmt.Errorf("http.Post failed: %v", err)
	}
	resData := swiftcResponseData{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		// TODO: send result
		return fmt.Errorf("json.Decode failed: %v", err)
	}
	if resData.Status != "success" {
		err := p.q.CreateSubmissionResult(ctx, db.CreateSubmissionResultParams{
			SubmissionID: int32(payload.SubmissionID),
			Status:       "compile_error",
			Stdout:       resData.Stdout,
			Stderr:       resData.Stderr,
		})
		if err != nil {
			// TODO: send result
			return fmt.Errorf("CreateTestcaseResult failed: %v", err)
		}
		p.results <- &TaskResultCompileSwiftToWasm{
			TaskPayload: &payload,
			Status:      "compile_error",
		}
		return fmt.Errorf("swiftc failed: %v", resData.Stderr)
	}

	// TODO: send result
	return nil
}

func (p *processor) processTaskCompileWasmToNativeExecutable(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadCompileWasmToNativeExecutable
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		// TODO: send result
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

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
	reqJson, err := json.Marshal(reqData)
	if err != nil {
		// TODO: send result
		return fmt.Errorf("json.Marshal failed: %v", err)
	}
	res, err := http.Post("http://worker:80/api/wasmc", "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		// TODO: send result
		return fmt.Errorf("http.Post failed: %v", err)
	}
	resData := wasmcResponseData{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		// TODO: send result
		return fmt.Errorf("json.Decode failed: %v", err)
	}
	if resData.Status != "success" {
		err := p.q.CreateSubmissionResult(ctx, db.CreateSubmissionResultParams{
			SubmissionID: int32(payload.SubmissionID),
			Status:       "compile_error",
			Stdout:       resData.Stdout,
			Stderr:       resData.Stderr,
		})
		if err != nil {
			// TODO: send result
			return fmt.Errorf("CreateTestcaseResult failed: %v", err)
		}
		p.results <- &TaskResultCompileWasmToNativeExecutable{
			TaskPayload: &payload,
			Status:      "compile_error",
		}
		return fmt.Errorf("wasmc failed: %v", resData.Stderr)
	}

	// TODO: send result
	return nil
}

func (p *processor) processTaskRunTestcase(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadRunTestcase
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		// TODO: send result
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	testcases, err := p.q.ListTestcasesByGameID(ctx, int32(payload.GameID()))
	if err != nil {
		// TODO: send result
		return fmt.Errorf("ListTestcasesByGameID failed: %v", err)
	}

	for _, testcase := range testcases {
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
			Stdin:       testcase.Stdin,
		}
		reqJson, err := json.Marshal(reqData)
		if err != nil {
			// TODO: send result
			return fmt.Errorf("json.Marshal failed: %v", err)
		}
		res, err := http.Post("http://worker:80/api/testrun", "application/json", bytes.NewBuffer(reqJson))
		if err != nil {
			// TODO: send result
			return fmt.Errorf("http.Post failed: %v", err)
		}
		resData := testrunResponseData{}
		if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
			// TODO: send result
			return fmt.Errorf("json.Decode failed: %v", err)
		}
		if resData.Status != "success" {
			err := p.q.CreateTestcaseResult(ctx, db.CreateTestcaseResultParams{
				SubmissionID: int32(payload.SubmissionID),
				TestcaseID:   testcase.TestcaseID,
				Status:       resData.Status,
				Stdout:       resData.Stdout,
				Stderr:       resData.Stderr,
			})
			if err != nil {
				// TODO: send result
				return fmt.Errorf("CreateTestcaseResult failed: %v", err)
			}
			p.results <- &TaskResultRunTestcase{
				TaskPayload: &payload,
				Status:      resData.Status,
			}
			return fmt.Errorf("testrun failed: %v", resData.Stderr)
		}
		if !isTestcaseResultCorrect(testcase.Stdout, resData.Stdout) {
			err := p.q.CreateTestcaseResult(ctx, db.CreateTestcaseResultParams{
				SubmissionID: int32(payload.SubmissionID),
				TestcaseID:   testcase.TestcaseID,
				Status:       "wrong_answer",
				Stdout:       resData.Stdout,
				Stderr:       resData.Stderr,
			})
			if err != nil {
				// TODO: send result
				return fmt.Errorf("CreateTestcaseResult failed: %v", err)
			}
			p.results <- &TaskResultRunTestcase{
				TaskPayload: &payload,
				Status:      "wrong_answer",
			}
			return fmt.Errorf("testrun failed: %v", resData.Stdout)
		}
	}

	p.results <- &TaskResultRunTestcase{
		TaskPayload: &payload,
		Status:      "success",
	}
	return nil
}

func isTestcaseResultCorrect(expectedStdout, actualStdout string) bool {
	expectedStdout = strings.TrimSpace(expectedStdout)
	actualStdout = strings.TrimSpace(actualStdout)
	return actualStdout == expectedStdout
}
