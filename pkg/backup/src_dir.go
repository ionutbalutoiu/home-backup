package backup

import (
	"context"
	"fmt"

	"github.com/ionutbalutoiu/home-backup/pkg/config"
)

func NewDirectorySourceBackup(params map[string]string) (SourceBackup, error) {
	dirParams := config.SrcDirectoryParams{}
	if err := dirParams.ParseParams(params); err != nil {
		return nil, fmt.Errorf("error parsing Directory source backup params: %w", err)
	}
	return &DirectorySourceBackup{Params: dirParams}, nil
}

type DirectorySourceBackup struct {
	Params config.SrcDirectoryParams
}

func (d *DirectorySourceBackup) Prepare(ctx context.Context) (string, error) {
	return d.Params.Path, nil
}

func (d *DirectorySourceBackup) Cleanup(ctx context.Context) error {
	return nil
}
