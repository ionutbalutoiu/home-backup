package backup

import "context"

// Source acquires an input for one backup run.
type Source interface {
	Open(context.Context) (Input, error)
}

// Input owns an acquired path and its release lifecycle.
type Input interface {
	Path() string
	Release(context.Context) error
}
