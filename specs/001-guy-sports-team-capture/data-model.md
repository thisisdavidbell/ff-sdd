# Data Model: Guy Sports Team Capture and Usage Analysis

## Core entities

### ManagerSnapshot
Represents the state of one manager's team at a specific capture time.

Fields:
- manager_id: string
- manager_name: string
- season: string
- captured_at: RFC3339 timestamp
- source: string
- players: []PlayerReference
- raw_url: string (optional)

Validation rules:
- manager_id must be stable across captures for the same manager.
- captured_at must be present and in UTC or RFC3339.
- players must be unique within a snapshot.
- source must be one of the supported upstream sources (initially Guy Sports).

### PlayerReference
A single footballer referenced in a manager's team.

Fields:
- player_id: string
- name: string
- team_name: string (optional)
- position: string (optional)

Validation rules:
- player_id should be stable enough to allow aggregation across snapshots.
- name and player_id should be kept as stored from the source to preserve traceability.

### TeamChangeEvent
A derived change between two snapshots for the same manager.

Fields:
- manager_id: string
- from_captured_at: RFC3339 timestamp
- to_captured_at: RFC3339 timestamp
- added_players: []PlayerReference
- removed_players: []PlayerReference
- change_count: int

Validation rules:
- added_players and removed_players are computed from the set difference between consecutive snapshots.
- only actual team differences generate a change event.
- unchanged snapshots are still retained but do not create a new change event.

### PlayerSelectionAggregate
Derived aggregate showing ownership by player across managers.

Fields:
- player_id: string
- player_name: string
- captured_at: RFC3339 timestamp
- manager_count: int
- managers: []string

Validation rules:
- manager_count equals the number of unique managers holding the player at that snapshot.
- history is built from the sequence of capture timestamps.

### ManagerChangeSummary
Derived summary of team changes for a manager.

Fields:
- manager_id: string
- total_changes: int
- latest_change_at: RFC3339 timestamp (optional)
- changed_since_last_snapshot: bool
- event_history: []TeamChangeEvent

Validation rules:
- total_changes is cumulative across all captured snapshots for the season.
- latest_change_at is blank if no changes have occurred.

## Storage model

### Raw data
- data/raw/<source>/<manager-id>/<timestamp>.yaml
- or a single YAML file per manager if the dataset remains manageable

The raw directory is append-only. Every new capture creates another snapshot, never overwriting the previous one.

### Processed data
- data/processed/player-ownership/<timestamp>.yaml
- data/processed/manager-changes/<timestamp>.yaml

Processed files are derived outputs; they may be regenerated and retained or discarded during local testing.

### Rendered output
- site/index.html or docs/index.html
- generated from processed data only, without live queries

The rendered HTML is the reader-facing artifact that should remain available as the last known-good version.

## Relationships
- One manager may have many ManagerSnapshot records across time.
- One snapshot contains many PlayerReference records.
- A TeamChangeEvent derives from two sequential snapshots for the same manager.
- Many PlayerSelectionAggregate records form the historical trend for a player.
- One ManagerChangeSummary summarizes all change events for a given manager.
