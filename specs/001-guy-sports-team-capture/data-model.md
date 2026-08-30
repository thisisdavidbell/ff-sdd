# Data Model: Guy Sports Team Capture and Usage Analysis

## Configuration Entity

### ProjectConfig (`config.yaml`)
Stores runtime configuration settings.

Fields:
- `season`: string (e.g. `"2026-27"`)
- `pages`: int (number of season table pages to capture, default `3`)
- `source_url`: string (base URL for the Guy Sports season table)

## Core entities

### ManagerSnapshot
Represents the state of one manager's team at a specific capture time.

Fields:
- `manager_id`: string
- `manager_name`: string (full text with spaces replaced by `_`)
- `team_name`: string (full text with spaces replaced by `_`)
- `captured_at`: RFC3339 timestamp
- `source`: string
- `players`: []PlayerReference
- `raw_url`: string (optional)

Validation rules:
- `manager_id` must be stable across captures for the same manager.
- `manager_name` and `team_name` must have spaces replaced with underscores (`_`).
- `captured_at` must be present and in UTC or RFC3339.
- `players` must be unique within a snapshot.
- `source` must be one of the supported upstream sources (initially Guy Sports).
- Stored under `data/<season>/raw/<team_name>_<manager_id>/<timestamp>.yaml`.
- Duplicate snapshots for a manager whose team is unchanged since the previous snapshot are skipped.

### PlayerReference
A single footballer referenced in a manager's team.

Fields:
- `player_id`: string
- `name`: string
- `position`: string
- `team_name`: string

Validation rules:
- `player_id` should be stable enough to allow aggregation across snapshots.
- `name`, `position`, and `team_name` should be retained to support human inspection of changes and to make manager change detail readable.

### TeamChangeEvent
A derived change between two snapshots for the same manager.

Fields:
- `manager_id`: string
- `from_captured_at`: RFC3339 timestamp
- `to_captured_at`: RFC3339 timestamp
- `added_players`: []PlayerReference
- `removed_players`: []PlayerReference
- `change_count`: int

Validation rules:
- `added_players` and `removed_players` are computed from the set difference between consecutive snapshots, and MUST include player `name` in addition to `player_id`.
- `change_count` represents the number of substitutions/transfers made in the event and MUST be calculated as `max(len(added_players), len(removed_players))` (e.g. 1 player added and 1 player removed equals 1 change, 2 added and 2 removed equals 2 changes).
- only actual team differences generate a change event.
- unchanged snapshots do not create a new change event.

### PlayerOwnershipFile (`data/<season>/processed/player-ownership.yaml`)
A single, unified human-readable YAML document containing current ownership counts and historical count trends across all capture timestamps for each player.

Fields:
- `season`: string
- `generated_at`: RFC3339 timestamp
- `players`: []PlayerOwnershipDetail
  - `player_id`: string
  - `player_name`: string
  - `current_count`: int
  - `history`: []PlayerOwnershipHistoryEntry
    - `captured_at`: RFC3339 timestamp
    - `manager_count`: int

Validation rules:
- At any capture timestamp $T$, a manager's active team composition is defined as their most recent snapshot captured on or before timestamp $T$.
- For every capture timestamp $T$, `manager_count` in `history` MUST be evaluated across the active teams of all known managers in the league as of timestamp $T$.
- `current_count` is derived from the latest recorded snapshot of each manager across the league, ensuring `current_count` equals the count at the final `captured_at` timestamp in `history`.
- `current_count`, the HTML player ownership table, and the final point in the historical trend chart MUST be 100% consistent with each other. Partial capture runs (where unchanged teams produce no new snapshot file) MUST NOT cause player counts to drop for players on unchanged teams.

### ManagerChangeSummary
Derived summary of team changes for a manager.

Fields:
- `manager_id`: string
- `manager_name`: string (spaces replaced by `_`)
- `team_name`: string (spaces replaced by `_`)
- `total_changes`: int
- `latest_change_at`: RFC3339 timestamp (optional)
- `changed_since_last_snapshot`: bool
- `event_history`: []TeamChangeEvent

Validation rules:
- `total_changes` is the cumulative sum of `change_count` for all `TeamChangeEvent` records in `event_history` (`total_changes = sum(event.change_count)`).
- `latest_change_at` is blank if no changes have occurred.
- Stored as `data/<season>/processed/manager-changes/<team_name>_<manager_id>.yaml`.

## Storage model

### Raw and processed data per season
- `config.yaml` (root configuration for season, pages, URLs)
- `data/<season>/raw/<team_name>_<manager_id>/<timestamp>.yaml`
- `data/<season>/processed/player-ownership.yaml` (single human-readable file with current counts and historical trend)
- `data/<season>/processed/manager-changes/<team_name>_<manager_id>.yaml`

The raw directory is append-only. Only distinct team states are added; unchanged rosters do not produce duplicate files. Processed files are derived outputs; the latest processed data currently present in the local repository state is the source used to generate the HTML render.

### Rendered output
- `docs/index.html`
- Generated from `player-ownership.yaml` and `manager-changes/*.yaml` without live queries.
- Player ownership table is sorted by manager count descending (highest first).
- Manager changes table displays Manager Name and Team Name, and MUST provide an interactive collapsible detail view per manager (expandable on click) that displays event timestamps, added player names (`+ Player Name`), removed player names (`- Player Name`), and event change counts for all events in `event_history`.
- Historical trends section is rendered as a large line chart showing the change in player counts over time across snapshot capture dates, featuring consistent (uniform/equidistant) horizontal gaps between dates on the x-axis and generous vertical height (e.g. 600px–800px) so many players can be comfortably distinguished.

## Relationships
- One manager may have many ManagerSnapshot records across time.
- One snapshot contains many PlayerReference records.
- A TeamChangeEvent derives from two sequential snapshots for the same manager.
- Many PlayerSelectionAggregate records form the historical trend for a player.
- One ManagerChangeSummary summarizes all change events for a given manager.
