package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type ExternalCommand struct {
	Command      []string
	CWD          string
	ReturnOutput bool
}

// ExecCommand executes an external command based on the provided ExternalCommand struct.
func ExecCommand(ctx context.Context, extCmd ExternalCommand) (string, error) {
	cmd := exec.CommandContext(ctx, extCmd.Command[0], extCmd.Command[1:]...)
	if extCmd.CWD != "" {
		cmd.Dir = extCmd.CWD
	}
	log.Debug(strings.Join(extCmd.Command, " "))
	if extCmd.ReturnOutput {
		output, err := cmd.CombinedOutput()
		return string(output), err
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	return "", err
}

// GetExitCode retrieves the exit code from an error returned by exec.Command.
func GetExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exitError.ExitCode()
	}
	// Non-exit error (e.g., command not found)
	return -1
}

// MountDevice mounts the device at the given path and returns the mount point.
func MountDevice(ctx context.Context, devicePath string) (string, error) {
	log.Debugf("detecting filesystem type for device %s", devicePath)
	output, err := ExecCommand(ctx, ExternalCommand{
		Command:      []string{"blkid", "-o", "value", "-s", "TYPE", devicePath},
		ReturnOutput: true,
	})
	if err != nil {
		return "", fmt.Errorf("detecting filesystem type for %s: %w", devicePath, err)
	}
	fsType := strings.TrimSpace(output)
	if fsType == "" {
		return "", fmt.Errorf("no filesystem detected on device %s", devicePath)
	}
	log.Debugf("creating temporary directory for mounting device %s", devicePath)
	tmpDir, err := os.MkdirTemp("", "lvm-backup_*")
	if err != nil {
		return "", err
	}
	// XFS requires "nouuid" to mount snapshots of the same filesystem
	var data string
	if fsType == "xfs" {
		data = "nouuid"
	}
	log.Debugf("mounting device %s (fstype=%s)", devicePath, fsType)
	if err := syscall.Mount(devicePath, tmpDir, fsType, syscall.MS_RDONLY, data); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("mounting %s on %s: %w", devicePath, tmpDir, err)
	}
	return tmpDir, nil
}

// UnmountDevice unmounts the filesystem at the given path.
func UnmountDevice(path string) error {
	log.Debugf("unmounting: %s", path)
	return syscall.Unmount(path, 0)
}

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
