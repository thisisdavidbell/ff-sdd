# AGENTS.md

## Repository Purpose

This repository provides a Go-based fantasy football data pipeline:

- `cmd/capture-guysports` gathers manager team snapshots into raw, append-only YAML data for GuySports.
- `cmd/process` derives player ownership and manager-change data.
- `cmd/render` produces the static HTML report in `docs/`.
- `run.sh` executes capture, processing, and rendering sequentially.

## Development Process

See [DEVELOPMENT-PROCESS.md](DEVELOPMENT-PROCESS.md) for the project's spec-driven development process.