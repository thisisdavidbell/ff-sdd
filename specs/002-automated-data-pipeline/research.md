# Research: Automated Data Pipeline & Capture Source Separation

**Feature**: [002-automated-data-pipeline](../spec.md)
**Date**: 2026-08-30

## Architectural Decisions

### 1. Source-Specific Capture Command Architecture

- **Decision**: Rename `cmd/capture` to `cmd/capture-guysports`.
- **Rationale**: Isolates provider-specific HTTP fetching and HTML parsing logic for GuySports. Moving the directory `cmd/capture` -> `cmd/capture-guysports` explicitly separates GuySports collection without requiring internal code changes in `main.go`. Allows future capture sources (e.g. `cmd/capture-dreamteam`) to be added cleanly.
- **Alternatives Considered**:
  - *Keep `cmd/capture` and use a `--source` flag*: Rejected because future sources will require different flags, parser configs, and endpoints. Source-specific entrypoints isolate dependencies and CLI contracts.
  - *Use Go subcommands (e.g., `capture guysports`)*: Rejected because the current repo pattern uses distinct binary targets in `bin/` (`bin/process`, `bin/render`).

### 2. Pipeline Execution Script (`run.sh`)

- **Decision**: Create a root-level executable shell script `run.sh` with strict error flags (`set -euo pipefail`).
- **Rationale**: Ensures portability between local environments and CI/CD runners. `set -e` ensures any failing step (build, capture, process, render) immediately halts the script, satisfying FR-002 and SC-002 (0% corrupted reports).
- **Alternatives Considered**:
  - *Makefile*: Rejected as `run.sh` is directly requested by the spec/user prompt and avoids Makefile tab/formatting quirks across environments.
  - *Go orchestrator binary*: Extra overhead without benefit for a simple 3-stage linear pipeline.

### 3. GitHub Actions Scheduled Workflow & Persistence

- **Decision**: Define `.github/workflows/schedule-run-action.yml` with:
  - Triggers: `schedule` with cron `0 18 * * 5` (Friday at 18:00 UTC) and `workflow_dispatch` (manual trigger).
  - Environment: `ubuntu-latest` with Go 1.22.x.
  - Execution: Checkout repo -> setup Go -> run `./run.sh` -> commit (`Automated weekly pipeline run`) and push modified files (`data/`, `docs/`) back to the branch using `GITHUB_TOKEN` git identity.
- **Rationale**: Meets FR-005, FR-006, and FR-008. Automated git commit preserves historical snapshots additively in accordance with Constitution Principle VI.
- **Alternatives Considered**:
  - *Commit via third-party Action*: Using standard `git` commands with `GITHUB_TOKEN` inside the workflow step keeps dependencies minimal and fully transparent.
