# Command Contracts

## Configuration

Configuration is loaded from `config.yaml` located in the repository root:
- `season`: active season identifier (e.g. `2026-27`)
- `pages`: number of season table pages to capture (e.g. `3`)
- `source_url`: base URL for the Guy Sports season table

## Capture command

Command: ./bin/capture --source guysports

Behavior:
- Reads configuration from `config.yaml` to determine season and number of pages (1 through N).
- Fetches each configured page from the season table, discovering and deduplicating managers.
- Formats manager and team names with spaces replaced by underscores (`_`).
- Checks each manager's dedicated directory `data/<season>/raw/<team_name>_<manager_id>/` for their latest snapshot.
- Skips writing a new snapshot if the manager's team is unchanged from their latest snapshot.
- Writes new timestamped manager snapshot YAML files under `data/<season>/raw/<team_name>_<manager_id>/<timestamp>.yaml` if the roster has changed or is new.
- Returns a summary of managers captured and files written.
- This is the only command in the normal workflow that may query the live Guy Sports site.

Failure modes:
- Upstream page unavailable: log the error and exit non-zero.
- Partial page failure: retain successful snapshots and record the failed source page.

## Process command

Command: ./bin/process

Behavior:
- Reads all raw YAML snapshots from manager directories under `data/<season>/raw/`.
- Computes latest player ownership counts and chronological count trends.
- Writes a single unified player ownership file: `data/<season>/processed/player-ownership.yaml`.
- Computes manager change summaries and writes individual files: `data/<season>/processed/manager-changes/<team_name>_<manager_id>.yaml`.
- Does not mutate raw capture history.
- Must not query Guy Sports or any other live source.

Failure modes:
- Empty or malformed input: warn, skip invalid records, and continue where possible.
- Missing capture data: exit non-zero with clear diagnostics.

## Render command

Command: ./bin/render

Behavior:
- Reads `data/<season>/processed/player-ownership.yaml` and manager changes from `data/<season>/processed/manager-changes/*.yaml`.
- Orders the Player Ownership table by manager count descending (highest count first), with ties broken alphabetically by player name.
- Renders the Manager Changes table with columns for Manager Name, Team Name, Total Changes, and Latest Change.
- Renders historical player trend data from the unified ownership file as a large line chart with consistent x-axis spacing between capture dates and generous vertical height to comfortably see many players.
- Writes static HTML output to `docs/index.html`.
- Uses embedded templates and static assets with no live API dependency.
- Must not query Guy Sports or any live source.

Failure modes:
- Missing processed data: fail clearly and instruct the user to run capture/process first.
- Template generation error: exit non-zero and preserve the previous good render.

## Reset/local cleanup policy

Local test runs may discard uncommitted generated working files, but must not modify the raw capture archive unless explicitly approved as a deliberate cleanup action.
