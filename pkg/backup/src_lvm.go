package backup

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ionutbalutoiu/home-backup/internal/utils"
	"github.com/ionutbalutoiu/home-backup/pkg/config"

	log "github.com/sirupsen/logrus"
)

const (
	defaultLVMSnapshotSize = "10G"
)

func NewLVMSourceBackup(params map[string]string) (SourceBackup, error) {
	lvmParams := config.SrcLVMParams{}
	if err := lvmParams.ParseParams(params); err != nil {
		return nil, fmt.Errorf("error parsing LVM source backup params: %w", err)
	}
	return &LVMSourceBackup{Params: lvmParams}, nil
}

type LVMSourceBackup struct {
	Params    config.SrcLVMParams
	mountPath string
}

func (l *LVMSourceBackup) Prepare(ctx context.Context) (retPath string, retErr error) {
	log.Debug("syncing filesystem before creating LVM snapshot")
	if _, err := utils.ExecCommand(ctx, utils.ExternalCommand{Command: []string{"sync"}}); err != nil {
		return "", fmt.Errorf("error syncing filesystem: %w", err)
	}
	log.Debug("creating LVM snapshot")
	cmd := utils.ExternalCommand{
		Command: []string{
			"lvcreate", "--snapshot",
			"--size", defaultLVMSnapshotSize,
			"--name", l.getSnapshotName(),
			l.getLVPath(),
		},
	}
	if _, err := utils.ExecCommand(ctx, cmd); err != nil {
		return "", fmt.Errorf("error creating LVM snapshot: %w", err)
	}
	defer func() {
		if retErr != nil {
			l.removeSnapshot(ctx)
		}
	}()
	log.Debug("mounting LVM snapshot")
	mountPath, err := utils.MountDevice(ctx, l.getSnapshotPath())
	if err != nil {
		return "", fmt.Errorf("error mounting LVM snapshot: %w", err)
	}
	l.mountPath = mountPath
	return l.mountPath, nil
}

func (l *LVMSourceBackup) Cleanup(ctx context.Context) error {
	var errs []error
	log.Debug("unmounting LVM snapshot")
	if err := utils.UnmountDevice(l.mountPath); err != nil {
		errs = append(errs, fmt.Errorf("unmounting snapshot at %s: %w", l.mountPath, err))
	}
	log.Debug("removing LVM snapshot")
	if err := l.removeSnapshot(ctx); err != nil {
		errs = append(errs, fmt.Errorf("removing snapshot: %w", err))
	}
	if l.mountPath != "" {
		log.Debugf("removing temporary mount path: %s", l.mountPath)
		if err := os.RemoveAll(l.mountPath); err != nil {
			errs = append(errs, fmt.Errorf("removing mount path %s: %w", l.mountPath, err))
		}
	}
	return errors.Join(errs...)
}

func (l *LVMSourceBackup) removeSnapshot(ctx context.Context) error {
	cmd := utils.ExternalCommand{
		Command: []string{"lvremove", "--force", l.getSnapshotPath()},
	}
	if _, err := utils.ExecCommand(ctx, cmd); err != nil {
		return fmt.Errorf("error removing LVM snapshot %s: %w", l.getSnapshotPath(), err)
	}
	return nil
}

func (l *LVMSourceBackup) getLVPath() string {
	return fmt.Sprintf("/dev/%s/%s", l.Params.VGName, l.Params.LVName)
}

func (l *LVMSourceBackup) getSnapshotName() string {
	return fmt.Sprintf("%s_backup_snapshot", l.Params.LVName)
}

func (l *LVMSourceBackup) getSnapshotPath() string {
	return fmt.Sprintf("/dev/%s/%s", l.Params.VGName, l.getSnapshotName())
}
