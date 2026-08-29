# Quickstart Guide

## Prerequisites
- Go 1.22 or newer
- Access to the project root
- Local filesystem write access for data capture output

## Local workflow

Important: no fake or example data may ever be written into the project’s canonical runtime directories. Test fixtures must be written to temp directories or a dedicated test-fixture area; the live data directories are reserved for real capture output only.

1. Build the binaries:
   - go build -o ./bin/capture ./cmd/capture
   - go build -o ./bin/process ./cmd/process
   - go build -o ./bin/render ./cmd/render
2. Capture raw team data from the live Guy Sports site:
   - ./bin/capture --source guysports
   - This is the only phase that may call Guy Sports.
   - Run this command yourself or through an authorized scheduler. AI assistants must not directly call Guy Sports; they may investigate only the locally stored capture artifacts.
3. Review the captured YAML snapshots under `data/2026-27/raw/`
4. Process the snapshots into derived analytics using only the local raw data:
   - ./bin/process
   - No live Guy Sports calls are allowed in this step.
5. Render the reader-facing HTML summary from the local processed data:
   - ./bin/render
   - No live Guy Sports calls are allowed in this step.
6. Open the generated HTML file in a browser to review the report

## Season configuration
- The current default season is `2026-27`.
- This is configured centrally in the project config and can be changed for the next season by updating the one season variable rather than editing multiple runtime paths.
- The capture, process, and render commands must all read the same configured season path so they stay consistent.

## Validation scenarios
- Verify raw files are timestamped and incremental, never overwritten.
- Verify processed counts match the manager snapshots in the selected time window.
- Verify rendered HTML includes ownership counts and manager change summaries.
- Verify the page loads without live network access.

## Testing note
- Early local experimentation may still use `go run` for quick iteration.
- Final project usage should rely on compiled binaries so the workflow is reproducible and matches the release-ready process.

## Local reset behavior
- Local experimentation may discard uncommitted generated processed output or the newest HTML render without committing those changes.
- The authoritative raw capture history remains intact and should not be removed unless intentional archive cleanup is planned.
- The last known-good HTML file remains the reader-facing artifact, and the render step uses the processed data currently available in the repo state at the time it runs.
