package backup

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const releaseTimeout = 2 * time.Minute

// LocalJob runs a local source-to-destination backup.
type LocalJob struct {
	source      Source
	destination Destination
}

// NewLocalJob constructs a local backup job.
func NewLocalJob(source Source, destination Destination) *LocalJob {
	return &LocalJob{source: source, destination: destination}
}

// Run opens the source, backs it up, and releases the acquired input.
func (j *LocalJob) Run(ctx context.Context) (retErr error) {
	input, err := j.source.Open(ctx)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), releaseTimeout)
		defer cancel()
		if err := input.Release(releaseCtx); err != nil {
			retErr = errors.Join(retErr, fmt.Errorf("release source: %w", err))
		}
	}()

	if err := j.destination.Backup(ctx, input.Path()); err != nil {
		return fmt.Errorf("backup destination: %w", err)
	}
	return nil
}

var _ Job = (*LocalJob)(nil)
