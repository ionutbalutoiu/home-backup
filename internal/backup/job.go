// Package backup defines technology-independent backup workflows.
package backup

import "context"

// Job is one independently executable backup.
type Job interface {
	Run(context.Context) error
}
