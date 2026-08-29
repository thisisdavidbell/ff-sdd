# Data Model: Guy Sports Team Capture and Usage Analysis

## Core entities

### ManagerSnapshot
Represents the state of one manager's team at a specific capture time.

Fields:
- manager_id: string
- manager_name: string
- captured_at: RFC3339 timestamp
- source: string
- players: []PlayerReference
- raw_url: string (optional)

Validation rules:
- manager_id must be stable across captures for the same manager.
- captured_at must be present and in UTC or RFC3339.
- players must be unique within a snapshot.
- source must be one of the supported upstream sources (initially Guy Sports).
- Season-specific storage is done by directory structure, not by a field on the snapshot itself.

### PlayerReference
A single footballer referenced in a manager's team.

Fields:
- player_id: string
- name: string
- position: string
- team_name: string

Validation rules:
- player_id should be stable enough to allow aggregation across snapshots.
- name, position, and team_name should be retained to support human inspection of changes and to make manager change detail readable.

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

Validation rules:
- manager_count equals the number of unique managers holding the player at that snapshot.
- history is built from the sequence of capture timestamps.
- The project stores only the count; it does not retain a list of managers who own each player.

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

### Raw and processed data per season
- data/<season>/raw/<source>/<manager-id>/<timestamp>.yaml
- data/<season>/processed/player-ownership/<timestamp>.yaml
- data/<season>/processed/manager-changes/<timestamp>.yaml
- or a single YAML file per manager under the season raw directory if the dataset remains manageable

The raw directory is append-only. Every new capture creates another snapshot, never overwriting the previous one. Processed files are derived outputs; the latest committed processed data is the canonical source used to generate the latest HTML render, and only uncommitted local experimental changes may be discarded during testing.

### Rendered output
- docs/index.html
- generated from the latest committed processed data only, without live queries

The rendered HTML is the reader-facing artifact that should remain available as the last known-good version and must be regenerated from the latest committed processed output rather than being treated as disposable local state.

## Relationships
- One manager may have many ManagerSnapshot records across time.
- One snapshot contains many PlayerReference records.
- A TeamChangeEvent derives from two sequential snapshots for the same manager.
- Many PlayerSelectionAggregate records form the historical trend for a player.
- One ManagerChangeSummary summarizes all change events for a given manager.
