package taskqueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
)

type processorWrapper struct {
	impl    processor
	results chan TaskResult
}

func newProcessorWrapper(impl processor) *processorWrapper {
	return &processorWrapper{
		impl:    impl,
		results: make(chan TaskResult),
	}
}

func (p *processorWrapper) processTaskCompileSwiftToWasm(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadCompileSwiftToWasm
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		err := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		p.results <- &TaskResultCompileSwiftToWasm{Err: err}
		return err
	}

	result, err := p.impl.doProcessTaskCompileSwiftToWasm(ctx, &payload)
	if err != nil {
		retryCount, _ := asynq.GetRetryCount(ctx)
		maxRetry, _ := asynq.GetMaxRetry(ctx)
		isRecoverable := !errors.Is(err, asynq.SkipRetry) && retryCount < maxRetry
		if !isRecoverable {
			p.results <- &TaskResultCompileSwiftToWasm{Err: err}
		}
		return err
	}
	p.results <- result
	return nil
}

func (p *processorWrapper) processTaskCompileWasmToNativeExecutable(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadCompileWasmToNativeExecutable
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		err := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		p.results <- &TaskResultCompileWasmToNativeExecutable{Err: err}
		return err
	}

	result, err := p.impl.doProcessTaskCompileWasmToNativeExecutable(ctx, &payload)
	if err != nil {
		retryCount, _ := asynq.GetRetryCount(ctx)
		maxRetry, _ := asynq.GetMaxRetry(ctx)
		isRecoverable := !errors.Is(err, asynq.SkipRetry) && retryCount < maxRetry
		if !isRecoverable {
			p.results <- &TaskResultCompileWasmToNativeExecutable{Err: err}
		}
		return err
	}
	p.results <- result
	return nil
}

func (p *processorWrapper) processTaskCreateSubmissionRecord(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadCreateSubmissionRecord
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		err := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		p.results <- &TaskResultCreateSubmissionRecord{Err: err}
		return err
	}

	result, err := p.impl.doProcessTaskCreateSubmissionRecord(ctx, &payload)
	if err != nil {
		retryCount, _ := asynq.GetRetryCount(ctx)
		maxRetry, _ := asynq.GetMaxRetry(ctx)
		isRecoverable := !errors.Is(err, asynq.SkipRetry) && retryCount < maxRetry
		if !isRecoverable {
			p.results <- &TaskResultCreateSubmissionRecord{Err: err}
		}
		return err
	}
	p.results <- result
	return nil
}

func (p *processorWrapper) processTaskRunTestcase(ctx context.Context, t *asynq.Task) error {
	var payload TaskPayloadRunTestcase
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		err := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		p.results <- &TaskResultRunTestcase{Err: err}
		return err
	}

	result, err := p.impl.doProcessTaskRunTestcase(ctx, &payload)
	if err != nil {
		retryCount, _ := asynq.GetRetryCount(ctx)
		maxRetry, _ := asynq.GetMaxRetry(ctx)
		isRecoverable := !errors.Is(err, asynq.SkipRetry) && retryCount < maxRetry
		if !isRecoverable {
			p.results <- &TaskResultRunTestcase{Err: err}
		}
		return err
	}
	p.results <- result
	return nil
}
