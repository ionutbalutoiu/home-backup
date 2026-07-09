package backup

import "context"

// Destination transfers an acquired path to its backing store.
type Destination interface {
	Backup(context.Context, string) error
}
