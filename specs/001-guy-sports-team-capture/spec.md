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

1. **Given** a season table with multiple pages of manager teams, **When** the scheduled capture runs, **Then** the system checks each page from page 1 through the available final page, including pages that contain fewer teams, and records the team composition for every manager found.
2. **Given** a manager has changed players since the previous capture, **When** the next capture runs, **Then** the system preserves both the prior and current team states and records the change.
3. **Given** the same manager appears on multiple season pages or the page count changes over time, **When** the capture runs, **Then** the system avoids losing historical records and treats the new data as an additional point in time rather than overwriting the past.

---

### User Story 2 - Show how often each player is selected by managers (Priority: P1)

A manager or analyst wants to see how often each player appears in the league. The system must produce two distinct but related outputs from the stored snapshots:

1. A current-state ownership view that looks only at each manager's latest team snapshot and counts how many managers currently hold each player.
2. A historical trend view that keeps the player-count totals from each capture event across the season so the change in popularity over time can be shown in a graph or trend table.

These two views must be easy to compare in a reader-friendly format.

**Why this priority**: This is the primary analytical output described for the first release and provides useful insight for fans and analysts.

**Independent Test**: Can be fully tested by processing captured team snapshots and verifying that the current-state counts use only the newest team snapshot for each manager while the historical trend preserves the count at each recorded capture time.

**Acceptance Scenarios**:

1. **Given** a set of historical team snapshots, **When** the aggregation step runs, **Then** the system uses each manager's most recent team snapshot to calculate the current player ownership count for each player, ordered from highest count to lowest, and presents the current snapshot in a clear summary.
2. **Given** a reader wants to understand player popularity over time, **When** they open the historical trend view, **Then** the system shows the count for each player across the stored capture timestamps in a clear, browsable format so changes over time can be reviewed.

---

### User Story 3 - Show each manager's team-change activity across the season (Priority: P1)

A manager wants to understand how many team changes they have made during the season. The system must track changes across captured snapshots and present each manager's change count in a readable summary, including when their latest change happened and whether any change has occurred since the previous snapshot.

**Why this priority**: This is a core decision-support metric described in the initial feature and directly supports the business value of the project.

**Independent Test**: Can be fully tested by comparing two consecutive manager snapshots and verifying that the system counts only actual team differences while preserving the season total for the manager.

**Acceptance Scenarios**:

1. **Given** a manager has a team snapshot and a later snapshot with different players, **When** the change analysis runs, **Then** the system counts the number of changes made since the previous snapshot.
2. **Given** a manager changes their team more than once across several captures, **When** the seasonal summary is generated, **Then** the system sums all changes made by that manager from the start of the season through the latest recorded snapshot.
3. **Given** a manager has made team changes across multiple snapshots, **When** the detailed view is opened, **Then** the system shows the date of the latest change, indicates whether a manager has changed since the last snapshot, and allows the reader to inspect each individual change event, including when it happened and what changed.

---

### User Story 4 - Present results in a GitHub-friendly and local-friendly format (Priority: P2)

An analyst needs a way to review computed results in a lightweight, shareable format without needing a live application. The processing step must produce the structured data needed for presentation, and the presentation layer must render the stored data in a professional, human-readable format that works both locally and in a GitHub Pages-style environment.

**Why this priority**: This enables the project to deliver value even before a full user-facing web experience exists and gives the project a clean, professional presentation without requiring live data access.

**Independent Test**: Can be fully tested by generating a rendered output and confirming that it is clear in a local browser or a static GitHub Pages-style host without any live network dependency.

**Acceptance Scenarios**:

1. **Given** processed data for player counts and change totals, **When** the presentation step runs, **Then** the system outputs a readable HTML-based summary that can be viewed locally and in a static GitHub Pages-style environment.
2. **Given** an HTML report is rendered for review, **When** it is opened locally or on a static host, **Then** it shows the latest player ownership counts ordered highest first, the historical count trend for each player, and the manager change summary with latest change dates and drill-down detail, and it never requests or depends on live data access.
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
- **FR-013**: The system MUST prefer a human-browsable storage model in which each manager has a dedicated data record or directory, with timestamped snapshots retained in chronological order for easy manual inspection.
- **FR-014**: The system MUST allow the raw historical data to be organized as either a single YAML file per manager or a per-manager directory containing timestamped snapshot files, provided the chosen structure remains efficient for processing and manual browsing.
- **FR-015**: The system MUST support processing of captured data into the structured outputs required for presentation, including latest player ownership counts, historical trend data for each player, and manager-level change summaries with timestamps.
- **FR-016**: The system MUST provide a reader-facing presentation that uses HTML-based rendering for a professional, static output that can be viewed locally and in a GitHub Pages-style environment without live data access.
- **FR-017**: The system MUST separate data capture, data processing, and presentation concerns into distinct commands or tools.
- **FR-018**: The system MUST support manual execution of capture commands for initial setup and testing before automation is enabled.
- **FR-019**: The system MUST explicitly exclude player score data and price data from this initial phase and focus only on players listed in users' teams.
- **FR-020**: The system MUST preserve historical records even when a team remains unchanged between two capture runs, so the record of time and state remains auditable.
- **FR-021**: The system MUST make it possible to clear local experimental capture and processed data without affecting the official stored dataset, so local experimentation can be safely reset while production data remains intact.
- **FR-022**: The system MUST use the latest processed data currently available in the local repository state as the input for the HTML render, whether that state is the latest committed data or the latest data produced by a local capture/process run, and the generated HTML render MUST remain a reader-facing artifact that is valuable to maintain and review.
- **FR-023**: The system MUST allow local testing to discard only uncommitted experimental changes to processed data or render outputs, while preserving the current repo state used for the latest render and ensuring the official raw historical capture data remains append-only and unmodified.
- **FR-024**: The system MUST keep all normal validation, unit tests, integration tests, and code-change checks offline and MUST NOT access the live Guy Sports or DreamTeamFC sites during routine validation.
- **FR-025**: The system MAY perform a limited live smoke test against the live site only when strictly necessary, such as after changing the retrieval logic, and only for the smallest possible validation: confirm that the connection to the live source works and that a small subset of the data is retrieved in the expected format, with Guy Sports as the only live site in scope for the initial implementation. This smoke test is not part of the normal validation cycle and is run only when the retrieval code has changed or when it has been explicitly requested.
- **FR-026**: The system MUST be designed to support future expansion for additional data sources and additional player-related datasets without breaking the historical capture model.

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
- **SC-003**: Player selection counts are reported accurately for the current state and selected historical windows, with no double-counting of a manager's team membership, and the change in counts over time is viewable from stored historical totals, ordered from highest to lowest in the latest summary.
- **SC-004**: Each manager's total team-change count is calculated accurately from sequential snapshots without overcounting unchanged states, and the total includes all change events since the start of the season.
- **SC-005**: Each manager's latest change date is clearly displayed, managers with changes since the last snapshot are highlighted, and a deeper drill-down view shows the change event details, including the date and the players added or removed.
- **SC-006**: Processed outputs are generated in a professional, repo-friendly HTML format that can be reviewed locally and in a GitHub Pages-style environment without a live data dependency, the latest valid render remains available to readers, and local test runs can be reset cleanly without disturbing the official captured dataset or the last known-good working render.
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
