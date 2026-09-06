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

`specs/` is an OKF v0.2 knowledge bundle. Every non-reserved Markdown artifact there MUST use the required OKF frontmatter profile defined in [DEVELOPMENT-PROCESS.md](DEVELOPMENT-PROCESS.md); `index.md` and `log.md` follow their reserved OKF roles. Do not extend this requirement to Markdown outside `specs/`.

## SDD Approval

After preparing the applicable SDD artifact, stop and obtain explicit user approval before changing implementation, configuration, automation, generated output, `ARCHITECTURE.md`, or any unrelated SDD artifact. Keep the artifact at `status: draft` while awaiting review; change it to `status: stable` only after approval. A stated user exception applies only to that current request and must be recorded in its SDD artifact.

An approved specification remains active until its scoped work is complete. Related user-requested refinements may be added to that active specification; completed historical specifications must not be amended to describe later work.

## AI Assistance

- The AI Agent should follow the users instructions at all times.
- The AI Agent can make suggestions and provide feedback.
- The AI Agent must follow the SDD Approval Gate and must not infer approval from silence, continued conversation, or a request to create, update, inspect, or discuss an SDD artifact.