#!/bin/sh
set -euo pipefail

TESTDATA_DIR="$1"

echo "--- Setting up test data ---"
mkdir -p /test/source-data/subdir
echo "hello world" > /test/source-data/file1.txt
echo "nested file content" > /test/source-data/subdir/file2.txt

echo "--- Running backup ---"
home-backup -config "${TESTDATA_DIR}/config.yaml" -log-level debug

echo "--- Verifying snapshot count ---"
SNAPSHOT_COUNT=$(restic --repo /tmp/restic-repo snapshots --json | grep -o '"short_id"' | wc -l)
if [ "${SNAPSHOT_COUNT}" -ne 1 ]; then
    echo "FAIL: expected 1 snapshot, got ${SNAPSHOT_COUNT}"
    exit 1
fi
echo "Snapshot count OK: ${SNAPSHOT_COUNT}"

echo "--- Verifying snapshot contents ---"
SNAPSHOT_FILES=$(restic --repo /tmp/restic-repo ls latest)
for EXPECTED in "/file1.txt" "/subdir/file2.txt"; do
    if ! echo "${SNAPSHOT_FILES}" | grep -q "${EXPECTED}"; then
        echo "FAIL: expected file not found in snapshot: ${EXPECTED}"
        exit 1
    fi
    echo "Found: ${EXPECTED}"
done
echo "Snapshot contents OK"
