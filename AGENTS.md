# AGENTS.md

## Repository Purpose

This repository provides a Go-based fantasy football data pipeline:

- `cmd/capture-guysports` gathers manager team snapshots into raw, append-only YAML data for GuySports.
- `cmd/process` derives player ownership and manager-change data.
- `cmd/render` produces the static HTML report in `docs/` for local viewing and github pages.
- `run.sh` executes capture, processing, and rendering sequentially.

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for the current application capabilities, high-level data flow, repository boundaries, and operating constraints.

## Development Process

See [DEVELOPMENT-PROCESS.md](DEVELOPMENT-PROCESS.md) for the project's spec-driven development process.

## SDD Routing

Before creating or changing SDD artifacts, read and follow the mandatory SDD Selection Gate in [DEVELOPMENT-PROCESS.md](DEVELOPMENT-PROCESS.md). Do not infer Full SpecKit SDD merely because its skills are available.

## AI Assistance

- The AI Agent should follow the users instructions at all times.
- The AI Agent can make suggestions and provide feedback.
- The AI Agent should never implement any changes without confirming with the user first.
- The AI should never perform its own tasks, or implement ideas without confirmation.