# Quickstart Guide

## Prerequisites
- Go 1.22 or newer
- Access to the project root
- Local filesystem write access for data capture output

## Local workflow

1. Build the binaries:
   - go build -o ./bin/capture ./cmd/capture
   - go build -o ./bin/process ./cmd/process
   - go build -o ./bin/render ./cmd/render
2. Capture raw team data:
   - ./bin/capture --source guysports
3. Review the captured YAML snapshots under data/raw/
4. Process the snapshots into derived analytics:
   - ./bin/process
5. Render the reader-facing HTML summary:
   - ./bin/render
6. Open the generated HTML file in a browser to review the report

## Validation scenarios
- Verify raw files are timestamped and incremental, never overwritten.
- Verify processed counts match the manager snapshots in the selected time window.
- Verify rendered HTML includes ownership counts and manager change summaries.
- Verify the page loads without live network access.

## Testing note
- Early local experimentation may still use `go run` for quick iteration.
- Final project usage should rely on compiled binaries so the workflow is reproducible and matches the release-ready process.

## Local reset behavior
- Local experimentation may discard generated processed output or the newest HTML render.
- The authoritative raw capture history remains intact and should not be removed unless intentional archive cleanup is planned.
- The last known-good HTML file remains the reader-facing artifact.
