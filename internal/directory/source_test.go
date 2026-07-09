package directory

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestSourceOpen(t *testing.T) {
	dir := t.TempDir()
	source := NewSource(dir)

	input, err := source.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if input.Path() != dir {
		t.Fatalf("Path() = %q, want %q", input.Path(), dir)
	}
	if err := input.Release(context.Background()); err != nil {
		t.Fatalf("Release() error = %v", err)
	}
}

func TestSourceOpenRejectsInvalidPaths(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "source-file")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "missing", path: file.Name() + "-missing", want: "stat directory source"},
		{name: "regular file", path: file.Name(), want: "is not a directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSource(tt.path).Open(context.Background())
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Open() error = %v, want substring %q", err, tt.want)
			}
		})
	}
}

func TestSourceOpenHonorsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewSource(t.TempDir()).Open(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Open() error = %v, want context.Canceled", err)
	}
}
