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

type ExecProcessor struct {
	q       *db.Queries
	results chan TaskExecResult
}

func NewExecProcessor(q *db.Queries, results chan TaskExecResult) *ExecProcessor {
	return &ExecProcessor{
		q:       q,
		results: results,
	}
}

func (p *ExecProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload TaskExecPlayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// TODO: upsert
	// Create submission record.
	submissionID, err := p.q.CreateSubmission(ctx, db.CreateSubmissionParams{
		GameID:   int32(payload.GameID),
		UserID:   int32(payload.UserID),
		Code:     payload.Code,
		CodeSize: int32(len(payload.Code)), // TODO: exclude whitespaces.
	})
	if err != nil {
		return fmt.Errorf("CreateSubmission failed: %v", err)
	}

	{
		type swiftcRequestData struct {
			MaxDuration int    `json:"max_duration_ms"`
			Code        string `json:"code"`
		}
		type swiftcResponseData struct {
			Result string `json:"result"`
			Stdout string `json:"stdout"`
			Stderr string `json:"stderr"`
		}
		reqData := swiftcRequestData{
			MaxDuration: 5000,
			Code:        payload.Code,
		}
		reqJson, err := json.Marshal(reqData)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %v", err)
		}
		res, err := http.Post("http://worker:80/api/swiftc", "application/json", bytes.NewBuffer(reqJson))
		if err != nil {
			return fmt.Errorf("http.Post failed: %v", err)
		}
		resData := swiftcResponseData{}
		if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
			return fmt.Errorf("json.Decode failed: %v", err)
		}
		if resData.Result != "success" {
			err := p.q.CreateTestcaseExecution(ctx, db.CreateTestcaseExecutionParams{
				SubmissionID: submissionID,
				TestcaseID:   nil,
				Status:       "compile_error",
				Stdout:       resData.Stdout,
				Stderr:       resData.Stderr,
			})
			if err != nil {
				return fmt.Errorf("CreateTestcaseExecution failed: %v", err)
			}
			p.results <- TaskExecResult{
				Task:   &payload,
				Result: "compile_error",
			}
			return fmt.Errorf("swiftc failed: %v", resData.Stderr)
		}
	}
	{
		type wasmcRequestData struct {
			MaxDuration int    `json:"max_duration_ms"`
			Code        string `json:"code"`
		}
		type wasmcResponseData struct {
			Result string `json:"result"`
			Stdout string `json:"stdout"`
			Stderr string `json:"stderr"`
		}
		reqData := wasmcRequestData{
			MaxDuration: 5000,
			Code:        payload.Code,
		}
		reqJson, err := json.Marshal(reqData)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %v", err)
		}
		res, err := http.Post("http://worker:80/api/wasmc", "application/json", bytes.NewBuffer(reqJson))
		if err != nil {
			return fmt.Errorf("http.Post failed: %v", err)
		}
		resData := wasmcResponseData{}
		if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
			return fmt.Errorf("json.Decode failed: %v", err)
		}
		if resData.Result != "success" {
			err := p.q.CreateTestcaseExecution(ctx, db.CreateTestcaseExecutionParams{
				SubmissionID: submissionID,
				TestcaseID:   nil,
				Status:       "compile_error",
				Stdout:       resData.Stdout,
				Stderr:       resData.Stderr,
			})
			if err != nil {
				return fmt.Errorf("CreateTestcaseExecution failed: %v", err)
			}
			p.results <- TaskExecResult{
				Task:   &payload,
				Result: "compile_error",
			}
			return fmt.Errorf("wasmc failed: %v", resData.Stderr)
		}
	}

	testcases, err := p.q.ListTestcasesByGameID(ctx, int32(payload.GameID))
	if err != nil {
		return fmt.Errorf("ListTestcasesByGameID failed: %v", err)
	}

	for _, testcase := range testcases {
		type testrunRequestData struct {
			MaxDuration int    `json:"max_duration_ms"`
			Code        string `json:"code"`
			Stdin       string `json:"stdin"`
		}
		type testrunResponseData struct {
			Result string `json:"result"`
			Stdout string `json:"stdout"`
			Stderr string `json:"stderr"`
		}
		reqData := testrunRequestData{
			MaxDuration: 5000,
			Code:        payload.Code,
			Stdin:       testcase.Stdin,
		}
		reqJson, err := json.Marshal(reqData)
		if err != nil {
			return fmt.Errorf("json.Marshal failed: %v", err)
		}
		res, err := http.Post("http://worker:80/api/testrun", "application/json", bytes.NewBuffer(reqJson))
		if err != nil {
			return fmt.Errorf("http.Post failed: %v", err)
		}
		resData := testrunResponseData{}
		if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
			return fmt.Errorf("json.Decode failed: %v", err)
		}
		if resData.Result != "success" {
			err := p.q.CreateTestcaseExecution(ctx, db.CreateTestcaseExecutionParams{
				SubmissionID: submissionID,
				TestcaseID:   testcase.TestcaseID,
				Status:       resData.Result,
				Stdout:       resData.Stdout,
				Stderr:       resData.Stderr,
			})
			if err != nil {
				return fmt.Errorf("CreateTestcaseExecution failed: %v", err)
			}
			p.results <- TaskExecResult{
				Task:   &payload,
				Result: resData.Result,
			}
			return fmt.Errorf("testrun failed: %v", resData.Stderr)
		}
		if !isTestcaseExecutionCorrect(testcase.Stdout, resData.Stdout) {
			err := p.q.CreateTestcaseExecution(ctx, db.CreateTestcaseExecutionParams{
				SubmissionID: submissionID,
				TestcaseID:   testcase.TestcaseID,
				Status:       "wrong_answer",
				Stdout:       resData.Stdout,
				Stderr:       resData.Stderr,
			})
			if err != nil {
				return fmt.Errorf("CreateTestcaseExecution failed: %v", err)
			}
			p.results <- TaskExecResult{
				Task:   &payload,
				Result: "wrong_answer",
			}
			return fmt.Errorf("testrun failed: %v", resData.Stdout)
		}
	}

	p.results <- TaskExecResult{
		Task:   &payload,
		Result: "success",
	}
	return nil
}

func isTestcaseExecutionCorrect(expectedStdout, actualStdout string) bool {
	expectedStdout = strings.TrimSpace(expectedStdout)
	actualStdout = strings.TrimSpace(actualStdout)
	return actualStdout == expectedStdout
}
