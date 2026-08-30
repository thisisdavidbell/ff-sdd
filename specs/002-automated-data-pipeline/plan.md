# Implementation Plan: Automated Data Pipeline & Capture Source Separation

**Branch**: `002-automated-data-pipeline` | **Date**: 2026-08-30 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-automated-data-pipeline/spec.md`

## Summary

Automate the Fantasy Football data collection pipeline by creating a unified pipeline execution script (`run.sh`), separating the GuySports capture CLI binary target (`cmd/capture-guysports`), and configuring a scheduled and manual GitHub Actions workflow (`.github/workflows/schedule-run-action.yml`). Updates to raw data, processed state, and rendered HTML reports will be committed and pushed back to the repository branch by the GitHub Actions workflow.

## Technical Context

**Language/Version**: Go 1.22.2, Bash shell, GitHub Actions YAML

**Primary Dependencies**: `gopkg.in/yaml.v3` v3.0.1, Go standard library (`net/http`, `os`, `path/filepath`, `flag`, `testing`)

**Storage**: Local raw YAML files (`data/2026-27/raw/`), processed YAML files (`data/2026-27/processed/`), static HTML (`docs/index.html`)

**Testing**: Go standard testing (`go test ./...`), integration tests in `tests/integration/` and unit tests in `tests/unit/`. **Development Directive**: At no point during the implementation of this feature should real live capture, process, or render pipeline commands be run automatically by agents/tooling against production endpoints/files. The user should be asked to manually run and verify live execution scenarios when desired and tidy up any unwanted changes after doing so.

**Target Platform**: Linux (GitHub Actions `ubuntu-latest`), macOS / Linux local development environments

**Project Type**: Data pipeline CLI & static site workflow

**Performance Goals**: Rapid sequential pipeline execution (capture, process, render)

**Constraints**: Additive historical data preservation per Constitution Principle VI; strict error handling (`set -euo pipefail`); commit and push changes back to repository branch in CI workflow step; no automated live production runs during feature implementation.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

1. **Spec-First Delivery**: PASS - Specification (`spec.md`), research (`research.md`), data model (`data-model.md`), contracts (`contracts/cli-commands.md`), and quickstart guide (`quickstart.md`) created before implementation.
2. **Data Integrity & Traceability**: PASS - `run.sh` enforces deterministic sequential execution: capture -> process -> render.
3. **Test-First Validation**: PASS - Tests will validate `capture-guysports` CLI behavior and script error handling.
4. **Observable, Reproducible Execution**: PASS - Script `run.sh` provides clear execution logs and non-zero failure signals.
5. **Simplicity and Change Control**: PASS - Simple bash execution script and dedicated CLI targets without extra framework overhead.
6. **Historical Data Preservation**: PASS - Capture step preserves historical raw snapshots additively in `data/2026-27/raw/`. CI workflow commits additions back to git.

## Project Structure

### Documentation (this feature)

```text
specs/002-automated-data-pipeline/
├── plan.md              # Implementation plan
├── research.md          # Technical decisions & research findings
├── data-model.md        # Data flow & repository artifact persistence model
├── quickstart.md        # Runnable validation guide
└── contracts/
    └── cli-commands.md  # CLI, script, and GitHub Actions contracts
```

### Source Code (repository layout)

```text
.github/
└── workflows/
    └── schedule-run-action.yml # Scheduled & manual workflow

bin/
├── capture-guysports           # Compiled GuySports capture binary
├── process                     # Compiled processing binary
└── render                      # Compiled rendering binary

cmd/
├── capture-guysports/
│   └── main.go                 # GuySports specific capture entrypoint
├── process/
│   └── main.go                 # Processing entrypoint
└── render/
    └── main.go                 # Rendering entrypoint

run.sh                          # Root pipeline execution script

data/                           # Versioned raw & processed YAML data
docs/                           # Rendered static HTML report
```

**Structure Decision**: Standard Go repository structure with root-level `run.sh` script and `.github/workflows/` runner configuration. `cmd/capture` renamed to `cmd/capture-guysports` for provider isolation.

## Complexity Tracking

*No constitution violations. Table left empty.*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
