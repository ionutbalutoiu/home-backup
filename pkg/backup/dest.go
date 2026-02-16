package backup

import (
	"context"
	"fmt"

	"github.com/ionutbalutoiu/home-backup/pkg/config"
)

type DestinationBackup interface {
	Create(ctx context.Context, backupPath string) error
}

func NewDestinationBackup(params map[string]string) (DestinationBackup, error) {
	if _, ok := params["type"]; !ok {
		return nil, fmt.Errorf("missing destination backup 'type' parameter")
	}
	switch params["type"] {
	case config.TypeRestic:
		return NewResticDestBackup(params)
	default:
		return nil, fmt.Errorf("unsupported destination backup type: %s", params["type"])
	}
}
