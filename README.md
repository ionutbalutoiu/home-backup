# home-backup

`home-backup` is a Linux CLI that runs sequential filesystem backups from a
typed YAML configuration. It currently supports directory and LVM snapshot
sources, with Restic as the destination.

## Usage

Build and run the command with a configuration file:

```sh
GOOS=linux go build -o ./build/home-backup ./cmd/home-backup
RESTIC_PASSWORD='...' ./build/home-backup -config ./config.yaml
```

Set `-log-level` to `debug`, `info`, `warn`, or `error` as needed. Directory
backups can run as an ordinary user; LVM backups require root and the Linux LVM
tools.

## Configuration

Each entry connects exactly one typed source to one typed destination:

```yaml
backups:
  - source:
      type: directory
      path: /srv/home
    destination:
      type: restic
      repo: /mnt/backups/home
      keep_last: 10
      group_by: host

  - source:
      type: lvm
      vg_name: vg0
      lv_name: home
    destination:
      type: restic
      repo: rclone:remote:home
      keep_last: 7
```

`keep_last` defaults to `10`, and `group_by` defaults to `host`. Configuration
decoding is strict: unknown fields, unsupported adapter types, invalid values,
and multiple YAML documents are rejected before any backup starts.

The Restic password and backend credentials are supplied through Restic's
normal environment variables or configuration. If a repository does not yet
exist, `home-backup` initializes it before creating the first snapshot.

## Architecture

The application uses explicit compile-time adapters, not runtime plugins:

- `internal/config` owns strict YAML decoding, defaults, and validation.
- `internal/app` is the composition root that maps typed configuration to
  concrete sources and destinations.
- `internal/backup` owns the small workflow contracts and sequential engine.
- `internal/directory` and `internal/lvm` implement source adapters.
- `internal/restic` implements the destination adapter.
- `internal/command` is the only external-command execution boundary.
- `cmd/home-backup` owns only process signals, standard streams, and exit code.

Each job opens a source input, sends its path to the destination, and releases
the input with a bounded cleanup context. LVM acquisition and cleanup state is
kept on the per-run input, so source objects remain reusable and stateless.

Adding a source or destination is a code change: define its typed config,
implement the narrow `backup.Source` or `backup.Destination` contract, and add
one explicit case in `internal/app`. This keeps supported integrations visible
to the compiler and avoids a registry or plugin framework in an application-only
module.

## Verification

Run the full suite on Linux:

```sh
go test ./...
go vet ./...
```

On a non-Linux host, run those commands in a Linux Go container. The repository
also includes `testdata/run-integration-test.sh` for the containerized Restic
smoke test.
