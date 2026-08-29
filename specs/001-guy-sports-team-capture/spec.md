# Feature Specification: Guy Sports Team Capture and Usage Analysis

**Feature Branch**: `001-guy-sports-team-capture`

**Created**: 2026-08-29

**Status**: Draft

**Input**: User description: "The project will have commands/scripts which are scheduled to run periodically to retrieve data. Initially these can be run manually, and as part of testing, to pull and store the data. there will be one for guysports and another for dreamteamfc. this first project will focus on guysports. it will retrieve the current team data for every manager playing the guy sports game. these can be found in the season table pages, e.g. page 1 is https://www.guysports.co.uk/guysports/season.php?page=1 . there are buttons on the page for each page, though currently only 3 have pages of managers. We do not need to pull the player score data when pulling the team data, or the price data. only the players listed in their team. Changes to the team should be captured, such that all data is preserved about the players in the team over time. This data will initially be used to show: counts of how many managers have selected each player, i.e. perhaps 30/32 managers have selected Haarland. how many changes to their team a manager has made. Over a season each manager only get 6 changes, so we will track how many they have made. There will be future pieces of work to pull additional data about the players scores and prices from https://www.guysports.co.uk/guysports/players.php?code=GK but this initial piece of work will focus on the players in the managers' teams. ideally the data will be stored in yaml files, though if this becomes too large, the ai may suggest a more efficient approach. the useful processed data (initially how many managers have selected each player, and how many changes each player had made ) will be presented in such a way that they can be viewed in github and locally. md files with something like mermaid, or github pages seem like the most appropriate. Havnig separate commands/tools for retrieving data and processing and presenting the data seems like a sensible approach."

## Terminology

- **Manager**: a person entered into the Guy Sports competition who selects a team.
- **Player**: a footballer who a manager may select for their team.
- **Team**: a manager's selection of 11 players that follows the game rules and
  team constraints for the competition.

These terms are used consistently throughout this specification and the project
vocabulary.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Capture and retain each manager's current team over time (Priority: P1)

A league administrator or analyst needs to understand how team composition changes across a season without relying on live site queries. The system must periodically capture the current team rosters for all Guy Sports managers and preserve each captured version as historical data.

**Why this priority**: Historical team data is the foundation for all later analysis, such as player popularity, team changes, and season trends.

**Independent Test**: Can be fully tested by running the periodic capture process for a known set of managers and verifying that each captured team snapshot is stored as a distinct historical record.

**Acceptance Scenarios**:

1. **Given** a season table configured for 3 pages in `config.yaml`, **When** the capture runs, **Then** the system checks pages 1, 2, and 3, discovers all managers across those pages, formats manager and team names with underscores replacing spaces, and saves snapshots under `data/<season>/raw/<team_name>_<manager_id>/<timestamp>.yaml`.
2. **Given** a manager has changed players since the previous capture, **When** the next capture runs, **Then** the system preserves the prior snapshot and records the new snapshot in the manager's directory.
3. **Given** a manager has not changed players since their latest snapshot, **When** the capture runs, **Then** the system detects that the team is unchanged and does not write a duplicate snapshot.
4. **Given** the same manager appears on multiple season pages or the page count changes over time, **When** the capture runs, **Then** the system deduplicates manager entries across pages before detail retrieval.

---

### User Story 2 - Show how often each player is selected by managers (Priority: P1)

A manager or analyst wants to see how often each player appears in the league. The system must produce two distinct but related outputs from the stored snapshots:

1. A current-state ownership view that looks only at each manager's latest team snapshot and counts how many managers currently hold each player.
2. A historical trend view that keeps the player-count totals from each capture event across the season so the change in popularity over time is displayed as a large line chart showing player counts over time with consistent x-axis spacing between capture dates and sufficient vertical height to comfortably see and distinguish many players.

These two views must be easy to compare in a reader-friendly format.

**Why this priority**: This is the primary analytical output described for the first release and provides useful insight for fans and analysts.

**Independent Test**: Can be fully tested by processing captured team snapshots and verifying that the current-state counts use only the newest team snapshot for each manager while the historical trend preserves the count at each recorded capture time and renders a large line chart with consistent date spacing and clear multi-player visibility.

**Acceptance Scenarios**:

1. **Given** a set of historical team snapshots across manager directories, **When** the aggregation step runs, **Then** the system outputs a single `data/<season>/processed/player-ownership.yaml` file containing both the current player ownership counts and historical count trends across capture timestamps.
2. **Given** an HTML report is rendered, **When** the player ownership table is generated, **Then** players are ordered by manager count descending (highest first), with ties broken alphabetically by player name.
3. **Given** historical trend data across multiple capture events, **When** the HTML report is generated, **Then** the historical trends section renders as a large line chart with consistent gaps between snapshot dates on the x-axis and generous height to clearly show changes across many players.

---

### User Story 3 - Show each manager's team-change activity across the season (Priority: P1)

A manager wants to understand how many team changes they have made during the season. The system must track changes across captured snapshots and present each manager's change count in a readable summary, including when their latest change happened and whether any change has occurred since the previous snapshot.

**Why this priority**: This is a core decision-support metric described in the initial feature and directly supports the business value of the project.

**Independent Test**: Can be fully tested by comparing two consecutive manager snapshots and verifying that the system counts only actual team differences while preserving the season total for the manager.

**Acceptance Scenarios**:

1. **Given** a manager has a team snapshot and a later snapshot with different players, **When** the change analysis runs, **Then** the system counts the number of changes made since the previous snapshot and saves the summary under `data/<season>/processed/manager-changes/<team_name>_<manager_id>.yaml`.
2. **Given** a manager changes their team more than once across several captures, **When** the seasonal summary is generated, **Then** the system sums all changes made by that manager from the start of the season through the latest recorded snapshot.
3. **Given** a manager has made team changes across multiple snapshots, **When** the detailed view is opened, **Then** the system shows the date of the latest change, indicates whether a manager has changed since the last snapshot, shows the manager name and team name, and allows the reader to inspect each individual change event.

---

### User Story 4 - Present results in a GitHub-friendly and local-friendly format (Priority: P2)

An analyst needs a way to review computed results in a lightweight, shareable format without needing a live application. The processing step must produce the structured data needed for presentation, and the presentation layer must render the stored data in a professional, human-readable format that works both locally and in a GitHub Pages-style environment.

**Why this priority**: This enables the project to deliver value even before a full user-facing web experience exists and gives the project a clean, professional presentation without requiring live data access.

**Independent Test**: Can be fully tested by generating a rendered output and confirming that it is clear in a local browser or a static GitHub Pages-style host without any live network dependency, including validating that the historical trends section displays as a large line chart with consistent x-axis date spacing.

**Acceptance Scenarios**:

1. **Given** processed data in `player-ownership.yaml` and `manager-changes/*.yaml`, **When** the presentation step runs, **Then** the system outputs a readable HTML-based summary (`docs/index.html`) that can be viewed locally and in a static GitHub Pages-style environment.
2. **Given** an HTML report is rendered for review, **When** it is opened locally or on a static host, **Then** it shows the player ownership counts sorted highest first, a large line chart representing the historical count trend for each player over time with consistent gaps between dates on the x-axis and sufficient vertical height to comfortably differentiate many players, and the manager change summary displaying Manager Name and Team Name with latest change dates, and it never requests or depends on live data access.
3. **Given** a local test run produces newer processed data or a new HTML render, **When** the experiment is discarded, **Then** the last known-good committed render remains available for the reader while the local test artifact may be reset without affecting the authoritative captured dataset.

### Edge Cases

- What happens when a manager is absent from a season page during a scheduled capture?
- What happens if a manager’s team is unchanged between two captures, but a valid historical record must still exist?
- What happens when an existing stored snapshot is incomplete or malformed?
- How does the system handle team pages that include more managers than the previously known total?
- How does the system handle duplicate records from multiple capture runs for the same manager and time window?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST support capture of Guy Sports manager team data from the season table pages at any interval, including irregular or scheduled runs, and each capture MUST be timestamped.
- **FR-002**: The system MUST capture the current team roster for each manager appearing on the Guy Sports season pages.
- **FR-003**: The system MUST preserve each captured team snapshot as historical data rather than overwriting prior state, including during normal operation and during local testing runs.
- **FR-004**: The system MUST store team snapshots additively so that past versions remain available for analysis and comparison over time, and each new capture from Guy Sports or DreamTeamFC is treated as a new data capture event rather than a replacement of the existing raw dataset.
- **FR-005**: The system MUST avoid live queries to the original Guy Sports data source during results presentation or manager-facing queries.
- **FR-006**: The system MUST treat captured data as the source of truth for manager-facing reporting and historical analysis.
- **FR-007**: The system MUST record the players in each manager's team for each capture event, including enough metadata to identify the team state at that point in time.
- **FR-008**: The system MUST track changes between successive team snapshots so that manager-level change counts can be calculated accurately regardless of the time gap between captures.
- **FR-009**: The system MUST compute the number of managers who selected each player across the relevant captured snapshots and time window.
- **FR-010**: The system MUST capture selection counts as historical values for each capture event so that changes in player popularity over time can be shown, not just the latest count.
- **FR-011**: The system MUST compute each manager's total team changes within a season using historical team snapshots.
- **FR-012**: The system MUST show, in human-readable output, the date of each manager's latest team change and provide a deeper view of individual change events, including the timestamp and the players affected.
- **FR-013**: The system MUST store each manager's raw historical data in a dedicated directory named by team name and manager ID (e.g. `<team_name>_<manager_id>`), containing timestamped snapshot YAML files to ensure human navigability and uniqueness.
- **FR-014**: The system MUST skip creating and storing a new snapshot for a manager whose team composition has not changed since their most recent recorded snapshot.
- **FR-015**: The system MUST store processed player ownership counts and historical player-count trends in a single, human-readable YAML file (`data/<season>/processed/player-ownership.yaml`) containing both the latest counts and the chronological historical trend for each player, allowing the render phase direct access to historical changes.
- **FR-016**: The system MUST format team names and manager names in raw snapshots and models with spaces replaced by underscores (`_`), preserving the full name text.
- **FR-017**: The system MUST name processed manager change summary files based on the team name with the manager ID suffix (e.g. `<team_name>_<manager_id>.yaml`).
- **FR-018**: The system MUST order player ownership in the rendered HTML output by manager count descending (highest count first), with ties broken alphabetically by player name.
- **FR-019**: The system MUST display manager changes in the rendered HTML table with both Manager Name and Team Name, rather than only manager ID.
- **FR-020**: The system MUST read runtime configuration from a dedicated configuration file (`config.yaml`), including the active season (e.g. `2026-27`) and the number of season table pages to capture (default `3`).
- **FR-021**: The system MUST capture data across all configured season table pages (pages 1 through N as specified in `config.yaml`, initially pages 1, 2, and 3), deduplicating managers discovered across pages before fetching rosters.
- **FR-022**: The system MUST render the historical trends section in the HTML report as a large line chart showing the change in player manager counts over time. The chart MUST be tall enough (e.g. at least 600px–800px tall) to comfortably see and distinguish many players, and MUST have consistent (uniform/equidistant) horizontal gaps between snapshot dates on the x-axis, regardless of whether the raw snapshots were captured with regular or irregular time intervals between them.
- **FR-023**: The system MUST never generate data in the data directory or html in the docs directory during testing. only execution of the capture, process or render commands can do this.
- **FR-024**: The system MUST use the latest processed data currently available in the local repository state as the input for the HTML render, whether that state is the latest committed data or the latest data produced by a local capture/process run, and the generated HTML render MUST remain a reader-facing artifact that is valuable to maintain and review.
- **FR-026**: The system MUST keep all normal validation, unit tests, integration tests, and code-change checks offline and MUST NOT access the live Guy Sports or DreamTeamFC sites during routine validation.
- **FR-027**: The system MAY perform a limited live smoke test against the live site only when strictly necessary, such as after changing the retrieval logic, and only for the smallest possible validation, but even then it must request permission first: confirm that the connection to the live source works and that a small subset of the data is retrieved in the expected format, with Guy Sports as the only live site in scope for the initial implementation. This smoke test is not part of the normal validation cycle and is run only when the retrieval code has changed or when it has been explicitly requested. AI assistants MUST NOT directly access Guy Sports; live capture is performed only by an explicitly user-invoked or authorized scheduled capture command.
- **FR-028**: The system MUST be designed to support future expansion for additional data sources and additional player-related datasets without breaking the historical capture model.
- **FR-029**: The system MUST NOT write sample, placeholder, fixture, or test data into the canonical project data directories used for real capture output. Example data for tests or demos MUST be stored in a temporary directory or a dedicated test-fixture location, never under `data/` or `docs/`.
- **FR-030**: The system MUST use the season directory naming convention `2026-27` for the current season and MUST support easy change-over to future seasons via a single configuration variable.
- **FR-031**: The system MUST default to live capture behavior for the Guy Sports source: the capture command, when explicitly run by a user or authorized scheduler, is the only phase that may query Guy Sports, and it must retrieve current manager team data from the live season pages and store it as append-only historical snapshots under the configured season directory. AI assistants MUST NOT directly access Guy Sports. The processing and HTML rendering phases MUST use only the locally stored captured data and MUST NOT call Guy Sports or any live source during normal execution.

### Key Entities

- **Manager Team Snapshot**: A record of one manager's team state at a specific capture time, including the season context and the roster of selected players.
- **Player Selection**: A player that appears in a manager's team snapshot, along with the manager and time context that identifies the selection event.
- **Season**: The league season or time window within which team snapshots are compared and aggregated.
- **Manager Change Count**: A derived value showing how many times a manager has changed their team between snapshots during the season and the timestamp of the most recent change.
- **Player Selection Aggregate**: A processed summary showing how many managers selected each player across the stored dataset, ordered by latest count and including the change in that count over time.
- **Manager Change Summary**: A processed summary showing each manager's total change count, latest change timestamp, and whether they changed since the previous snapshot.
- **Team Change Event**: A specific change between two snapshots, including the timestamp, the manager, the players removed, and the players added.
- **Presentation Output**: A generated HTML-based presentation used to view processed results locally or in a GitHub Pages-style static host without live data access.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The system captures team snapshots for all known managers across all available Guy Sports season pages during each scheduled run, including pages with fewer than the maximum number of teams.
- **SC-002**: Historical records remain available for all processed snapshots, allowing analysts to compare team states across the season without losing prior data.
- **SC-003**: Player selection counts are reported accurately for the current state and selected historical windows, with no double-counting of a manager's team membership, and the change in counts over time is viewable from stored historical totals, ordered from highest to lowest in the latest summary and rendered as a large line chart with consistent gaps between snapshot dates on the x-axis and ample vertical height to comfortably distinguish many players.
- **SC-004**: Each manager's total team-change count is calculated accurately from sequential snapshots without overcounting unchanged states, and the total includes all change events since the start of the season.
- **SC-005**: Each manager's latest change date is clearly displayed, managers with changes since the last snapshot are highlighted, and a deeper drill-down view shows the change event details, including the date and the players added or removed.
- **SC-006**: Processed outputs are generated in a professional, repo-friendly HTML format that can be reviewed locally and in a GitHub Pages-style environment without a live data dependency, featuring a large, readable historical trends line chart with uniform x-axis date spacing, the latest valid render remains available to readers, and local test runs can be reset cleanly without disturbing the official captured dataset or the last known-good working render.
- **SC-007**: The initial release supports Guy Sports as the primary source while leaving room for future expansion to additional data sources and derived metrics.

## Assumptions

- The first release focuses on Guy Sports only; DreamTeamFC capture is planned as a future or parallel feature and not required for this v1 scope.
- The system may store data in structured, versioned files such as YAML if the dataset size remains manageable; a more efficient format can be adopted if scale requires it.
- The preferred raw storage approach is a YAML file per manager, or, if needed for scale, a per-manager directory with timestamped snapshot files for each capture event.
- The season table pages are the authoritative source for current user roster data during the initial implementation phase.
- Historical capture can be performed in scheduled runs or manual execution during testing, with the same processing path used in both modes, and all captures are timestamped regardless of interval.
- Capture frequency may be regular or irregular; processing and presentation must still correctly show the historical sequence of team data and change events.
- User queries and presentation views rely only on persisted historical data and not on direct access to the upstream website.
- The data collection process is expected to preserve both state and change history, even when managers do not change their teams between consecutive captures.
- Local testing may generate temporary data in a disposable local working set, and a reset command may be used to discard the local generated data while preserving the official captured dataset used for scheduled production updates.
- The project may eventually run the capture and processing tools via GitHub Actions or similar automation, but this is a future operational consideration rather than a requirement for the initial feature scope.
- Raw capture data is append-only and must never be overwritten by normal operation or local test runs; every later fetch from Guy Sports or DreamTeamFC creates a new historical capture record that can be preserved, deleted, or rolled back independently.
- Processed data and rendered HTML outputs are derived artifacts that are important to readers because they make the analysis accessible, but only uncommitted local experimental changes may be discarded; the latest processed data currently available in the repo state is the source used to generate the latest HTML report, while the official historical raw dataset remains untouched.
- Routine validation MUST use fixtures and offline stored data only. Live access to Guy Sports or DreamTeamFC is forbidden during standard test runs and code-change validation. Only a narrowly scoped live smoke test is permitted when retrieval logic has changed or when an explicit manual request indicates it is needed, and that smoke test must be limited to confirming connectivity and a very small subset of expected format.
