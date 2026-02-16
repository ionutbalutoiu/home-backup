# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

home-backup is a Go CLI tool that automates home server backup operations. It reads a YAML config file describing backup jobs (each with a source and destination) and executes them. It must run as root. The tool is containerized and published as a multi-arch Docker image (amd64/arm64).

## Platform

This is a Linux-only tool. All builds must use `GOOS=linux`.

## Build Commands

```bash
# Build the binary
GOOS=linux go build -o ./build/home-backup ./cmd/home-backup

# Build Docker image
docker build -t home-backup .

# Run tests (none exist yet, but standard Go test command)
GOOS=linux go test ./...

# Vet / format
GOOS=linux go vet ./...
gofmt -l .
```

## Architecture

The codebase follows an interface-based plugin pattern for backup sources and destinations:

- **`cmd/home-backup/main.go`** — Entry point: parses `-config` flag, loads YAML config, calls `backup.CreateBackups()`
- **`pkg/config/`** — YAML config loading. Source/destination params are parsed as `map[string]string` and validated into typed param structs (`SrcLVMParams`, `SrcDirectoryParams`, `DestResticParams`)
- **`pkg/backup/`** — Core backup logic with two interfaces:
  - `SourceBackup` (`Prepare() string`, `Cleanup() error`) — prepares a filesystem path to back up. Implementations: `LVMSourceBackup` (creates LVM snapshot, mounts it), `DirectorySourceBackup` (returns path as-is)
  - `DestinationBackup` (`Create(backupPath string) error`) — performs the backup. Implementation: `ResticDestinationBackup` (runs restic backup + prune)
  - Factory functions `NewSourceBackup()` / `NewDestinationBackup()` switch on the `type` parameter
- **`internal/utils/`** — Shared utilities: command execution, mount/unmount (via syscall), filesystem type detection (via blkid), file existence checks

**Backup flow:** For each job: source.Prepare() → destination.Create(path) → source.Cleanup(). LVM sources create a temporary snapshot, mount it read-only, back it up, then clean up the snapshot.

## Key Dependencies

- `github.com/sirupsen/logrus` — structured logging
- `gopkg.in/yaml.v3` — config parsing
- External binaries at runtime: `restic`, `rclone` (bundled in Docker image), LVM tools (`lvcreate`, `lvremove`)

## CI/CD

GitHub Actions workflow (`.github/workflows/home-backup.yaml`) triggers on semver git tags, builds a multi-arch Docker image, and pushes to GitHub Container Registry. Renovate (`renovate.json`) auto-updates Go modules, Dockerfile versions, and GHA versions.

## Adding New Source/Destination Types

1. Add a type constant in `pkg/config/constants.go`
2. Create a params struct in `pkg/config/` with validation
3. Implement the `SourceBackup` or `DestinationBackup` interface in `pkg/backup/`
4. Register it in the corresponding factory function (`NewSourceBackup` or `NewDestinationBackup`)
