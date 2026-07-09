package backup

import (
	"context"
	"errors"
	"fmt"
)

// Engine executes configured backup jobs sequentially.
type Engine struct {
	jobs []Job
}

// NewEngine constructs an engine from independent backup jobs.
func NewEngine(jobs ...Job) *Engine {
	return &Engine{jobs: append([]Job(nil), jobs...)}
}

// Run executes every job unless the context is cancelled.
func (e *Engine) Run(ctx context.Context) error {
	var errs []error
	for i, job := range e.jobs {
		if err := ctx.Err(); err != nil {
			errs = append(errs, fmt.Errorf("stop before backup %d: %w", i+1, err))
			break
		}
		if err := job.Run(ctx); err != nil {
			errs = append(errs, fmt.Errorf("backup %d: %w", i+1, err))
		}
	}
	return errors.Join(errs...)
}
