# Quickstart Guide

## Prerequisites
- Go 1.22 or newer
- Access to the project root
- Local filesystem write access for data capture output

## Local workflow

1. Capture raw team data:
   - go run ./cmd/capture --source guysports
2. Review the captured YAML snapshots under data/raw/
3. Process the snapshots into derived analytics:
   - go run ./cmd/process
4. Render the reader-facing HTML summary:
   - go run ./cmd/render
5. Open the generated HTML file in a browser or serve it locally:
   - python3 -m http.server 8000 -d ./site

## Validation scenarios
- Verify raw files are timestamped and incremental, never overwritten.
- Verify processed counts match the manager snapshots in the selected time window.
- Verify rendered HTML includes ownership counts and manager change summaries.
- Verify the page loads without live network access.

## Local reset behavior
- Local experimentation may discard generated processed output or the newest HTML render.
- The authoritative raw capture history remains intact and should not be removed unless intentional archive cleanup is planned.
- The last known-good HTML file remains the reader-facing artifact.
