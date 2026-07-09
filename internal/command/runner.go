package command

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/exec"
)

// Runner executes external commands using os/exec.
type Runner struct {
	logger *slog.Logger
}

// NewRunner constructs a command runner.
func NewRunner(logger *slog.Logger) *Runner { return &Runner{logger: logger} }

// Run executes a command and captures stdout, stderr, and exit status.
func (r *Runner) Run(ctx context.Context, spec Spec) (Result, error) {
	if spec.Name == "" {
		return Result{ExitCode: -1}, errors.New("command name is required")
	}
	if r.logger != nil {
		r.logger.DebugContext(ctx, "execute command",
			"name", spec.Name,
			"args", spec.Args,
			"dir", spec.Dir,
		)
	}

	cmd := exec.CommandContext(ctx, spec.Name, spec.Args...)
	cmd.Dir = spec.Dir
	if len(spec.Env) > 0 {
		cmd.Env = append(os.Environ(), spec.Env...)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := Result{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: 0}
	if err == nil {
		return result, nil
	}
	result.ExitCode = -1
	var processErr *exec.ExitError
	if errors.As(err, &processErr) {
		result.ExitCode = processErr.ExitCode()
	}
	cause := err
	if contextErr := ctx.Err(); contextErr != nil {
		cause = errors.Join(err, contextErr)
	}
	return result, &ExitError{Spec: spec, Result: result, Err: cause}
}
