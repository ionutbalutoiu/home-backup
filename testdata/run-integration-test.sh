#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

IMAGE_NAME="home-backup:integration-test"
CONTAINER_TESTDATA_DIR="/testdata"

echo "=== Building Docker image ==="
docker build -t "${IMAGE_NAME}" "${PROJECT_ROOT}"

echo "=== Running integration test ==="
docker run --rm \
    -e RESTIC_PASSWORD=integration-test-password \
    -v "${SCRIPT_DIR}:${CONTAINER_TESTDATA_DIR}:ro" \
    --entrypoint /bin/sh \
    "${IMAGE_NAME}" \
    "${CONTAINER_TESTDATA_DIR}/smoke-test/run.sh" "${CONTAINER_TESTDATA_DIR}/smoke-test"

echo "=== ALL TESTS PASSED ==="
