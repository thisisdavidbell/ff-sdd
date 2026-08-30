# Interface Contract: CLI Commands & Automation Script

**Feature**: [002-automated-data-pipeline](../spec.md)
**Date**: 2026-08-30

## 1. GuySports Data Capture Command (`capture-guysports`)

### Invocation
```bash
go run ./cmd/capture-guysports [flags]
# OR
./bin/capture-guysports [flags]
```

### Command Flags
- `-config string`: Path to configuration file (default: `config.yaml`).
- `-source string`: Data source identifier (default: `guysports`).
- `-output string`: Target raw directory path (default: derived path `data/<season>/raw` based on active season).
- `-pages int`: Number of season pages to crawl (default: evaluated from `config.yaml`).

### Behavior & Exit Codes
- **0**: Capture completed successfully (or skipped unchanged snapshots).
- **1**: Critical error (e.g. invalid config, network failure across all pages, unable to write snapshot).
- **2**: Usage / invalid flag error.

---

## 2. Pipeline Execution Script (`run.sh`)

### Invocation
```bash
./run.sh
```

### Behavior & Execution Sequence
1. Compiles binaries or executes commands in sequence:
   1. `capture-guysports`
   2. `process`
   3. `render`
2. Configured with `set -euo pipefail`.
3. Stops immediately on non-zero exit code from any stage.

### Exit Codes
- **0**: All three stages (`capture-guysports` -> `process` -> `render`) succeeded.
- **Non-zero**: Script failed at the first failing stage; subsequent stages were skipped.

---

## 3. GitHub Actions Workflow Contract (`.github/workflows/schedule-run-action.yml`)

### Triggers
```yaml
on:
  schedule:
    - cron: '0 18 * * 5' # Every Friday at 18:00 UTC
  workflow_dispatch:     # Manual triggering from GitHub UI / API
```

### Workflow Steps Contract
1. **Checkout Repository**: `actions/checkout@v4` with write permissions.
2. **Setup Go Environment**: `actions/setup-go@v5` with Go 1.22.x.
3. **Execute Pipeline**: Execute `./run.sh`.
4. **Commit & Push Changes**:
   - Detect changes in `data/` and `docs/`.
   - Commit message: `Automated weekly pipeline run`
   - Push to active repository branch.
