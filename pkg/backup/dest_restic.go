package backup

import (
	"context"
	"fmt"

	"github.com/ionutbalutoiu/home-backup/internal/utils"
	"github.com/ionutbalutoiu/home-backup/pkg/config"

	log "github.com/sirupsen/logrus"
)

const resticExitCodeRepoNotFound = 10

func NewResticDestBackup(params map[string]string) (DestinationBackup, error) {
	resticParams := config.DestResticParams{}
	if err := resticParams.ParseParams(params); err != nil {
		return nil, fmt.Errorf("error parsing Restic destination backup params: %w", err)
	}
	return &ResticDestinationBackup{Params: resticParams}, nil
}

type ResticDestinationBackup struct {
	Params config.DestResticParams
}

func (r *ResticDestinationBackup) Create(ctx context.Context, backupPath string) error {
	log.Debug("checking if restic repository exists")
	cmd := utils.ExternalCommand{
		Command:      []string{"restic", "--repo", r.Params.Repo, "cat", "config"},
		ReturnOutput: true,
	}
	if _, err := utils.ExecCommand(ctx, cmd); err != nil {
		if utils.GetExitCode(err) == resticExitCodeRepoNotFound {
			log.Debug("restic repository does not exist")
			log.Debug("initializing restic repository")
			cmd = utils.ExternalCommand{
				Command: []string{"restic", "--repo", r.Params.Repo, "init"},
			}
			if _, err := utils.ExecCommand(ctx, cmd); err != nil {
				return fmt.Errorf("failed to initialize restic repository: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check if restic repository exists: %w", err)
		}
	}
	log.Debug("creating restic backup")
	cmd = utils.ExternalCommand{
		Command: []string{"restic", "--repo", r.Params.Repo, "backup", "."},
		CWD:     backupPath,
	}
	if _, err := utils.ExecCommand(ctx, cmd); err != nil {
		return fmt.Errorf("failed to create restic backup: %w", err)
	}
	log.Debug("pruning old restic backups")
	cmd = utils.ExternalCommand{
		Command: []string{"restic", "--repo", r.Params.Repo, "forget", "--group-by", r.Params.GroupBy, "--keep-last", fmt.Sprintf("%d", r.Params.KeepLast), "--prune"},
	}
	if _, err := utils.ExecCommand(ctx, cmd); err != nil {
		return fmt.Errorf("failed to prune old restic backups: %w", err)
	}
	return nil
}
