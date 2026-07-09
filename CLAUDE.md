# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

home-backup is a Go CLI tool that automates home server backup operations. It reads a strict, typed YAML config describing source-to-destination jobs and executes them sequentially. Directory jobs can run unprivileged; LVM jobs require root. The tool is containerized and published as a multi-arch Docker image (amd64/arm64).

## Platform

This is a Linux-only tool. All builds must use `GOOS=linux`.

## Build Commands

```bash
# Build the binary
GOOS=linux go build -o ./build/home-backup ./cmd/home-backup

# Build Docker image
docker build -t home-backup .

# Run tests and static checks on Linux
go test ./...
go vet ./...

# Format
gofmt -l .
```

On a non-Linux host, run tests and vet inside a Linux Go container; setting
`GOOS=linux` locally cross-compiles test binaries but cannot execute them.

## Architecture

This is an application-only module. Sources and destinations are compile-time
extensions composed explicitly by the application, not third-party runtime
plugins.

- **`cmd/home-backup/main.go`** — Process lifecycle only: signal context, standard streams, and exit status.
- **`internal/app/`** — Composition root and CLI option parsing. This is the only package that selects concrete source and destination adapters.
- **`internal/config/`** — Strict YAML decoding into typed source and destination variants, plus defaults and structural validation. YAML nodes do not escape this package.
- **`internal/backup/`** — Workflow contracts (`Job`, `Source`, `Input`, and `Destination`), the local job lifecycle, and the sequential engine.
- **`internal/command/`** — Structured `os/exec` boundary with separate stdout/stderr capture and typed exit errors.
- **`internal/directory/`** — Directory source adapter.
- **`internal/lvm/`** — Linux-only LVM snapshot source and read-only system mounter.
- **`internal/restic/`** — Restic destination adapter and retention sequence.

**Backup flow:** `Source.Open` → `Destination.Backup(input.Path())` →
`Input.Release`. Cleanup uses an uncancelled, bounded context and joins cleanup
errors with backup errors. Per-run LVM state belongs to the returned input, not
the reusable source.

## Key Dependencies

- Go standard library `log/slog` — structured logging
- `gopkg.in/yaml.v3` — strict config parsing
- `golang.org/x/sys/unix` — Linux mount and unmount syscalls
- External binaries at runtime: `restic`, `rclone` (bundled in Docker image), LVM tools (`lvcreate`, `lvremove`)

## CI/CD

GitHub Actions workflow (`.github/workflows/home-backup.yaml`) triggers on semver git tags, builds a multi-arch Docker image, and pushes to GitHub Container Registry. Renovate (`renovate.json`) auto-updates Go modules, Dockerfile versions, and GHA versions.

## Adding New Source/Destination Types

1. Add the typed configuration variant and strict decoder fields in `internal/config`.
2. Implement `backup.Source` or `backup.Destination` in a focused adapter package.
3. Add the explicit construction case in `internal/app/wiring.go`.
4. Add decoder, adapter, wiring, and lifecycle tests. Do not introduce a registry, dependency-injection framework, or runtime plugin mechanism.
