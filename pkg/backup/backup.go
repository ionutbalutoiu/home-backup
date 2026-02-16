package backup

import (
	"context"
	"errors"

	"github.com/ionutbalutoiu/home-backup/pkg/config"

	log "github.com/sirupsen/logrus"
)

func CreateBackups(ctx context.Context, config *config.Config) error {
	var errs []error
	for _, backup := range config.Backups {
		log.Debugf("creating backup from %s to %s", backup.Source, backup.Destination)
		if err := createSingleBackup(ctx, backup); err != nil {
			log.Errorf("backup failed: %v", err)
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func createSingleBackup(ctx context.Context, backup config.Backup) error {
	log.Debug("setup source backup handler")
	srcBackup, err := NewSourceBackup(backup.Source)
	if err != nil {
		return err
	}
	log.Debug("setup destination backup handler")
	destBackup, err := NewDestinationBackup(backup.Destination)
	if err != nil {
		return err
	}
	log.Debug("setting up source backup (e.g., creating volume snapshot)")
	backupPath, err := srcBackup.Prepare(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := srcBackup.Cleanup(ctx); err != nil {
			log.Warnf("failed to clean up source backup %s: %v", backup.Source, err)
		}
	}()
	log.Debug("creating backup at destination")
	if err := destBackup.Create(ctx, backupPath); err != nil {
		return err
	}
	return nil
}
