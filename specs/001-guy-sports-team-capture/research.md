# Research: Guy Sports Team Capture and Usage Analysis

## Decision: Go is the default implementation language

We will implement the first release in Go. This aligns with the repository's data-processing nature, the need for repeatable CLI scripts, typed models, and straightforward HTML rendering without introducing a larger runtime dependency.

### Rationale
- Go gives reliable cross-platform execution on Linux/macOS development machines and CI.
- It is well-suited to scheduled or manual capture jobs that fetch source data, store timestamped snapshots, and process historical records.
- The standard library provides HTTP, filesystem, and templating support needed for static HTML output.
- The project is a data pipeline rather than a user-facing web app, so a small CLI-first tool is a good fit.

### Alternative considered
- Python was considered because of YAML parsing and quick scripting, but it introduces more cross-environment variance and less explicit type safety for historical data contracts.
- A non-Go implementation would only be approved if the team explicitly documents a concrete requirement that Go cannot satisfy and records the exception.

## Decision: YAML is the preferred raw storage format

Raw captured snapshots will be stored in YAML files, organized by manager and timestamp, with append-only historical records.

### Rationale
- YAML is easy for humans to inspect and diff in a Git repository.
- It supports structured snapshots for manager teams and change metadata without requiring a database.
- It matches the requirement that raw historical data be browsable and auditable.

### Alternative considered
- JSON would be machine-friendly but less readable in Git reviews.
- A database would improve scale but adds operational complexity and conflicts with the preferred human-browsable storage model.
- If YAML later proves inefficient at scale, a different storage format may be proposed, but only with explicit approval and documentation.

## Decision: HTML is the preferred reader-facing presentation

The generated presentation will be static HTML that can be opened locally and also deployed to GitHub Pages.

### Rationale
- It satisfies the requirement for a lightweight, shareable output that does not require a live backend.
- It works well for generated tables, charts, and summaries without introducing a runtime dependency.
- It is static-host friendly and repository-friendly.

### Alternatives considered
- Markdown with Mermaid was considered as a human-readable option but is less polished for a static reader experience and less ideal for GitHub Pages hosting.
- A live web app was rejected because the project explicitly prefers persisted historical data and static presentation without live data access.

## Decision: capture, processing, and rendering remain separate commands

The pipeline will be split into distinct tools or commands:
- capture: fetch raw team snapshots from source
- process: derive ownership counts and change summaries
- render: generate HTML static report

This keeps source fidelity, derived data, and presentation output cleanly separated.

## Working assumptions
- Every capture is timestamped and additive; no raw data overwrite.
- Local testing may regenerate processed and rendered output, but the last known-good render can be retained.
- The canonical data source remains persisted raw snapshots, not live upstream queries.
- The initial release focuses on Guy Sports; DreamTeamFC remains future work unless explicitly added.
