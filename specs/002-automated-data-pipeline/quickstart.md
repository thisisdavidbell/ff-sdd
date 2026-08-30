# Quickstart & Validation Guide: Automated Data Pipeline

**Feature**: [002-automated-data-pipeline](../spec.md)
**Date**: 2026-08-30

This guide describes how to validate the automated data pipeline and capture source separation locally and in CI/CD.

> **Important Verification Directive**:
> AI agents and automated tools MUST NOT execute live `capture-guysports`, `process`, or `render` commands against production endpoints/files during feature development.
> The user is requested to manually execute and verify these validation scenarios, and then clean up any unwanted working tree changes afterwards (e.g., via `git restore data/ docs/`).

## Prerequisites

- Go 1.22+ installed locally.
- Git repository workspace checked out on branch `002-automated-data-pipeline`.
- Internet connectivity (for GuySports capture tests) or mock environment.

---

## Manual Validation Scenarios (User Executed)

### Scenario 1: Verify Direct GuySports Capture Target

Validate that `capture-guysports` can be executed directly as a standalone CLI tool.

```bash
# Build binary target
go build -o bin/capture-guysports ./cmd/capture-guysports

# Run capture tool
./bin/capture-guysports -pages 1
```

**Expected Outcome**:
- Command completes with exit code `0`.
- Captured snapshots are stored under `data/2026-27/raw/` (or unchanged snapshots are reported as skipped).

---

### Scenario 2: Verify End-to-End Pipeline Script Execution (`run.sh`)

Validate that `run.sh` executes build/capture, process, and render stages sequentially.

```bash
# Make run.sh executable
chmod +x run.sh

# Run end-to-end pipeline
./run.sh
```

**Expected Outcome**:
- Output displays sequential progress through GuySports capture, data processing, and report rendering.
- Exit code is `0`.
- `docs/index.html` and `data/2026-27/processed/` files are updated with latest state.

---

### Scenario 3: Verify Pipeline Halt on Failure (Error Handling)

Validate that `run.sh` halts immediately if any stage fails.

```bash
# Test with an invalid flag/config to force capture failure
./bin/capture-guysports -source invalid_source || echo "Capture failed as expected"

# Verify run.sh halts when a step fails
# (e.g. by temporarily pointing config.yaml to an invalid URL or output path)
```

**Expected Outcome**:
- If `capture-guysports` exits with non-zero code, `process` and `render` steps are NOT executed.
- `run.sh` exits with non-zero exit code.

---

### Scenario 4: Validate GitHub Actions Workflow Syntax

Validate that the GitHub Actions workflow file syntax is valid.

```bash
# Verify workflow file exists
ls -la .github/workflows/schedule-run-action.yml

# Optional: Validate syntax using actionlint if available
actionlint .github/workflows/schedule-run-action.yml || true
```

**Expected Outcome**:
- `.github/workflows/schedule-run-action.yml` exists, sets schedule `0 18 * * 5`, enables `workflow_dispatch`, calls `./run.sh`, and handles git commit/push.
