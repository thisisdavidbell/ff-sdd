---
type: Quickstart
title: "Quickstart Guide"
description: Historical quickstart artifact for 001-guy-sports-team-capture.
tags: [sdd, quickstart, 001-guy-sports-team-capture]
status: stable
feature: 001-guy-sports-team-capture
sdd_approach: full-speckit
input_summary: Historical quickstart artifact.
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
---

# Quickstart Guide

## Prerequisites
- Go 1.22 or newer
- Access to the project root
- Local filesystem write access for data capture output

## Configuration (`config.yaml`)
Configuration is managed in `config.yaml` at the project root:
```yaml
season: "2026-27"
pages: 3
source_url: "https://www.guysports.co.uk/guysports/season.php"
```
- `season`: sets the directory path for raw snapshots and processed output (`data/<season>/...`).
- `pages`: specifies the number of season table pages to capture (pages 1 through N).

## Local workflow

Important: no fake or example data may ever be written into the project’s canonical runtime directories. Test fixtures must be written to temp directories or a dedicated test-fixture area; the live data directories are reserved for real capture output only.

1. Build the binaries:
   - go build -o ./bin/capture ./cmd/capture
   - go build -o ./bin/process ./cmd/process
   - go build -o ./bin/render ./cmd/render
2. Capture raw team data from the live Guy Sports site:
   - ./bin/capture --source guysports
   - This pulls data from pages 1..3 as configured in `config.yaml`.
   - Each manager's snapshots are placed in `data/2026-27/raw/<team_name>_<manager_id>/<timestamp>.yaml`.
   - If a manager's team is unchanged from their latest snapshot, duplicate snapshots are skipped.
   - Run this command yourself or through an authorized scheduler. AI assistants must not directly call Guy Sports; they may investigate only the locally stored capture artifacts.
3. Review the captured YAML snapshots under `data/2026-27/raw/`
4. Process the snapshots into derived analytics using only the local raw data:
   - ./bin/process
   - Generates a single `data/2026-27/processed/player-ownership.yaml` (with current counts and historical trend).
   - Generates `data/2026-27/processed/manager-changes/<team_name>_<manager_id>.yaml`.
   - No live Guy Sports calls are allowed in this step.
5. Render the reader-facing HTML summary from the local processed data:
   - ./bin/render
   - Renders `docs/index.html` with player ownership sorted highest count first, manager changes showing Manager Name and Team Name, and a large historical trends line chart with consistent x-axis spacing between capture dates.
   - No live Guy Sports calls are allowed in this step.
6. Open the generated HTML file in a browser to review the report

## Validation scenarios
- Verify raw files are stored under per-manager directories `data/2026-27/raw/<team_name>_<manager_id>/`.
- Verify unchanged teams do not produce redundant duplicate snapshots.
- Verify processed player ownership is consolidated into `data/2026-27/processed/player-ownership.yaml`.
- Verify rendered HTML orders players by count descending, displays manager and team names, and renders historical trends as a large line chart with consistent date spacing and comfortable height for distinguishing many players.
- Verify the page loads without live network access.

## Testing note
- Early local experimentation may still use `go run` for quick iteration.
- Final project usage should rely on compiled binaries so the workflow is reproducible and matches the release-ready process.

## Local reset behavior
- Local experimentation may discard uncommitted generated processed output or the newest HTML render without committing those changes.
- The authoritative raw capture history remains intact and should not be removed unless intentional archive cleanup is planned.
- The last known-good HTML file remains the reader-facing artifact, and the render step uses the processed data currently available in the repo state at the time it runs.
