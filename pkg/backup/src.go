package backup

import (
	"context"
	"fmt"

	"github.com/ionutbalutoiu/home-backup/pkg/config"
)

type SourceBackup interface {
	Prepare(ctx context.Context) (string, error)
	Cleanup(ctx context.Context) error
}

func NewSourceBackup(params map[string]string) (SourceBackup, error) {
	if _, ok := params["type"]; !ok {
		return nil, fmt.Errorf("missing source backup 'type' parameter")
	}
	switch params["type"] {
	case config.TypeLVM:
		return NewLVMSourceBackup(params)
	case config.TypeDirectory:
		return NewDirectorySourceBackup(params)
	default:
		return nil, fmt.Errorf("unsupported source backup type: %s", params["type"])
	}
}
