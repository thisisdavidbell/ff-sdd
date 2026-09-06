# ff-sdd
A Go-based fantasy football data pipeline for a GuySports league. It records manager team snapshots over time, calculates player ownership and team changes, and publishes a static report that makes league activity easy to inspect.

## What It Does

The pipeline has three stages:

- **Capture** downloads the configured GuySports league pages and stores append-only manager team snapshots as YAML.
- **Process** turns those snapshots into current and historical player ownership plus manager-change summaries.
- **Render** generates [docs/index.html](docs/index.html), a self-contained report with ownership movement, historical trends, expandable team-change details, and responsive section navigation.

Run all three stages with `./run.sh`, or run them individually when working on a particular stage.

## Usage

### Direct Stage Commands
- **Capture (GuySports)**: `go run ./cmd/capture-guysports`
- **Process**: `go run ./cmd/process`
- **Render**: `go run ./cmd/render`

### Automated End-to-End Pipeline
- **Run Pipeline**: `./run.sh`

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for the current capabilities, data flow, data locations, repository boundaries, configuration, and operating constraints.

## Development

This repository uses spec-driven development. See [DEVELOPMENT-PROCESS.md](DEVELOPMENT-PROCESS.md) for the workflow and [specs/index.md](specs/index.md) for its OKF v0.2 feature knowledge bundle. All non-reserved Markdown artifacts within `specs/` use the required OKF frontmatter profile.

