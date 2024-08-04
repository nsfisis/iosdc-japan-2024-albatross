package taskqueue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type ExecProcessor struct {
}

func (processor *ExecProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload TaskExecPlayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// TODO
	return nil
}

func NewExecProcessor() *ExecProcessor {
	return &ExecProcessor{}
}
