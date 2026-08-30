#!/usr/bin/env bash
set -euo pipefail

# Automated Fantasy Football Data Pipeline
# Runs capture, process, and render stages sequentially via `go run`.
# Exits immediately on any stage error.

echo "=== Starting Data Capture (GuySports) ==="
go run ./cmd/capture-guysports "$@"

echo "=== Processing Player Ownership & Manager Changes ==="
go run ./cmd/process "$@"

echo "=== Rendering Static HTML Report ==="
go run ./cmd/render "$@"

echo "=== Pipeline Completed Successfully ==="
