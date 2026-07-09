package backup

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type jobFunc func(context.Context) error

func (f jobFunc) Run(ctx context.Context) error { return f(ctx) }

func TestEngineContinuesAndJoinsErrors(t *testing.T) {
	firstErr := errors.New("first")
	secondRan := false
	engine := NewEngine(
		jobFunc(func(context.Context) error { return firstErr }),
		jobFunc(func(context.Context) error { secondRan = true; return nil }),
	)

	err := engine.Run(context.Background())
	if !errors.Is(err, firstErr) || !secondRan {
		t.Fatalf("Run() error = %v, secondRan = %v", err, secondRan)
	}
	if !strings.Contains(err.Error(), "backup 1") {
		t.Fatalf("Run() error = %v, want job index", err)
	}
}

func TestEngineStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	secondRan := false
	engine := NewEngine(
		jobFunc(func(context.Context) error { cancel(); return nil }),
		jobFunc(func(context.Context) error { secondRan = true; return nil }),
	)

	err := engine.Run(ctx)
	if !errors.Is(err, context.Canceled) || secondRan {
		t.Fatalf("Run() error = %v, secondRan = %v", err, secondRan)
	}
}
