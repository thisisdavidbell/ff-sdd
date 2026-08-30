# AGENTS.md

## Repository Purpose

This repository provides a Go-based fantasy football data pipeline:

- `capture` gathers manager team snapshots into raw, append-only YAML data.
- `process` derives player ownership and manager-change data.
- `render` produces the static HTML report in `docs/`.

## Development Process

This project uses SpecKit and follows spec-driven development (SDD).

Before implementing **any** change, including minor fixes or maintenance, update the relevant feature artifacts in `specs/<feature>/` first. Keep `spec.md`, `plan.md`, `tasks.md`, and related documents (such as contracts, data models, and research) accurate for the change.

Use the SpecKit workflow to specify, plan, task, and implement work. Do not make implementation-only changes that are absent from the relevant SpecKit Markdown artifacts.