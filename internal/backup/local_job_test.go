package backup

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type sourceFunc func(context.Context) (Input, error)

func (f sourceFunc) Open(ctx context.Context) (Input, error) { return f(ctx) }

type destinationFunc func(context.Context, string) error

func (f destinationFunc) Backup(ctx context.Context, path string) error { return f(ctx, path) }

type fakeInput struct {
	path    string
	release func(context.Context) error
}

func (f *fakeInput) Path() string { return f.path }

func (f *fakeInput) Release(ctx context.Context) error { return f.release(ctx) }

func TestLocalJobRunOrdersLifecycle(t *testing.T) {
	var calls []string
	job := NewLocalJob(
		sourceFunc(func(context.Context) (Input, error) {
			calls = append(calls, "open")
			return &fakeInput{path: "/snapshot", release: func(context.Context) error {
				calls = append(calls, "release")
				return nil
			}}, nil
		}),
		destinationFunc(func(_ context.Context, path string) error {
			calls = append(calls, "backup:"+path)
			return nil
		}),
	)

	if err := job.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	want := []string{"open", "backup:/snapshot", "release"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestLocalJobRunJoinsBackupAndReleaseErrors(t *testing.T) {
	backupErr := errors.New("backup failed")
	releaseErr := errors.New("release failed")
	job := NewLocalJob(
		sourceFunc(func(context.Context) (Input, error) {
			return &fakeInput{path: "/snapshot", release: func(context.Context) error { return releaseErr }}, nil
		}),
		destinationFunc(func(context.Context, string) error { return backupErr }),
	)

	err := job.Run(context.Background())
	if !errors.Is(err, backupErr) || !errors.Is(err, releaseErr) {
		t.Fatalf("Run() error = %v, want both errors", err)
	}
}

func TestLocalJobRunReleasesWithUncancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	job := NewLocalJob(
		sourceFunc(func(context.Context) (Input, error) {
			return &fakeInput{path: "/snapshot", release: func(ctx context.Context) error {
				if err := ctx.Err(); err != nil {
					t.Fatalf("release context error = %v", err)
				}
				return nil
			}}, nil
		}),
		destinationFunc(func(context.Context, string) error {
			cancel()
			return context.Canceled
		}),
	)

	if err := job.Run(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context cancellation", err)
	}
}
