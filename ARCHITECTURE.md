# Architecture

## Purpose

This repository is a Go-based fantasy football data pipeline for a GuySports league. It captures manager Team (roster) snapshots, derives player ownership and manager-change history, and publishes a static HTML report for local viewing or GitHub Pages.

## Agreed Terminology

- **Manager**: a person entered into the Guy Sports competition who selects a team.
- **Player**: a footballer who a manager may select for their team.
- **Team**: a manager's selection of 11 players that follows the competition rules and team constraints.

In reader-facing documentation, **Team** is the canonical term. The first relevant reference may include **roster** as an explanatory synonym. Code identifiers, data fields, and historical specifications may retain `roster` when changing them would add churn or break an existing contract.

## Current Capabilities

- Capture manager and team data from the configured GuySports season pages.
- Preserve team snapshots as append-only YAML, skipping unchanged teams.
- Calculate current and historical player ownership across managers.
- Detect team additions and removals between snapshots.
- Render ownership trends and team changes as an interactive, responsive static report with light and dark themes, including compact mobile chart sizing and a mobile Team changes view that omits the Manager column.
- Show each player's latest ownership direction and last change time, and prioritize managers with the most recent changes.
- Provide persistent grouped section navigation: a full-height desktop sidebar with Statistics, Trends, and Settings groups, and a closed-by-default hamburger menu positioned below the report title and timestamp on narrow screens.

## High-Level Architecture

The application is a sequential three-stage pipeline. Each stage is a separate Go command, and YAML files provide the durable hand-off between stages.

```text
-----------------------------------------------------------------------
capture-guysports  -->  data/<season>/raw/<team>_<manager-id>/*.yaml
                                      |
                                      v
process             -->  data/<season>/processed/player-ownership.yaml
                         data/<season>/processed/manager-changes/*.yaml
                                      |
                                      v
render              -->  docs/index.html
-----------------------------------------------------------------------
      |
      v
GuySports Report HTML
```

### Capture

`cmd/capture-guysports`:

- Reads the configured source pages.
- Extracts manager links and fetches each manager's team.
- Writes timestamped `models.ManagerSnapshot` YAML files.

Capture is source-specific today: the supported source is GuySports. Raw files are grouped by manager so the history remains append-only and can be processed repeatedly.

### Processing

`cmd/process`:

- Discovers raw snapshots from the capture phase through `internal/storage`.
- Uses `internal/processing` to reconstruct the latest active team for each manager at each capture time.
- Aggregates player ownership history.
- Derives manager-level team changes.
- Writes one consolidated ownership document and one change summary per manager as YAML.

### Rendering

`cmd/render`:

- Reads the processed YAML from the processing phase.
- Converts it into the report view model.
- Delegates HTML generation to `internal/render`.
- Retains UTC timestamps in capture and processed YAML, then converts reader-facing report timestamps to `Europe/London` at minute precision.
- Requires IANA timezone data for `Europe/London` and fails rendering rather than publishing timestamps in a fallback timezone when that data is unavailable.

The result is a self-contained static page containing:

- Ownership counts with their latest increase or decrease indicator and last change time, initially previewing ten rows with a reversible control to show all players.
- Historical trends.
- Team changes ordered by most recent activity, initially previewing ten manager rows with a reversible control to show all teams and expandable details.
- Browser-side interaction, theme persistence through the navigation Settings group, and active section navigation.

## Repository Boundaries

- `cmd/`: executable stage entry points and pipeline orchestration.
- `internal/config/`: configuration loading, environment overrides, and derived paths.
- `internal/models/`: YAML-serializable contracts shared across stages.
- `internal/capture/`: GuySports HTML parsing, team fetching, and snapshot helpers.
- `internal/storage/`: snapshot discovery and YAML file I/O.
- `internal/processing/`: ownership aggregation and change detection.
- `internal/render/`: static HTML generation.
- `data/<season>/`: raw inputs and processed intermediate outputs.
- `docs/`: generated published HTML web page.
- `specs/`: an OKF v0.2 knowledge bundle holding SDD feature concepts and supporting artifacts, with YAML frontmatter for discovery and validation.
- `tests/unit/` and `tests/integration/`: parser, model, processing, rendering, configuration, script, and persistence coverage.

## Configuration

`config.yaml` provides the default pipeline settings:

- `season`: the season-specific data directory, currently `2026-27`.
- `pages`: the number of GuySports season pages to capture, currently `3`.
- `source_url`: the GuySports season page used as the capture source.

Configuration can also be overridden with:

- `GUYSPORTS_CONFIG`: path to an alternate configuration file.
- `GUYSPORTS_SEASON`: season override used for derived data paths.

## Operating Model

Run stages individually:

- `go run ./cmd/capture-guysports`
- `go run ./cmd/process`
- `go run ./cmd/render`

Run `./run.sh` to execute the stages in order. The pipeline stops when a stage returns an error.

The GitHub Actions scheduled pipeline runs daily at 01:54 in the `Europe/London` timezone, avoiding the start-of-hour scheduling peak and retaining the same local time through daylight-saving changes.

Repository changes follow one of the SDD approaches defined in `DEVELOPMENT-PROCESS.md`. SDD artifacts are maintained as OKF concepts within `specs/`, with `index.md` and `log.md` supporting discovery and history. GitHub Copilot pull-request review checks that applicable SDD artifacts accompany implementation and behavior changes, and that this current-state architecture document is updated if needed to remain in sync with the code behaviour.

## Current Constraints

- GuySports is the only supported capture source, and capture depends on the source HTML structure and network availability.
- Raw snapshots are intentionally manager-scoped and append-only; processing assumes the latest known snapshot represents a manager when that manager has no new capture at a later time.
- The report is static output generated from processed YAML; its browser-side charts and controls are not a separate server application.
- Rendering depends on IANA timezone data for `Europe/London`; deployment environments without it cannot generate a report.
- Generated data and report output are season-specific, so the configured season and corresponding `data/<season>/` layout must agree.
