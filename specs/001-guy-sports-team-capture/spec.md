# Feature Specification: Guy Sports Team Capture and Usage Analysis

**Feature Branch**: `001-guy-sports-team-capture`

**Created**: 2026-08-29

**Status**: Draft

**Input**: User description: "The project will have commands/scripts which are scheduled to run periodically to retrieve data. Initially these can be run manually, and as part of testing, to pull and store the data. there will be one for guysports and another for dreamteamfc. this first project will focus on guysports. it will retrieve the current team data for every user playing the guy sports game. these can be found in the season table pages, e.g. page 1 is https://www.guysports.co.uk/guysports/season.php?page=1 . there are buttons on the page for each page, though currently only 3 have pages of users. We do not need to pull the player score data when pulling the team data, or the price data. only the players listed in their team. Changes to the team should be captured, such that all data is preserved about the players in the team over time. This data will initially be used to show: counts of how many users have selected each player, i.e. perhaps 30/32 users have selected Haarland. how many changes to their team a user has made. Over a season each user only get 6 changes, so we will track how many they have made. There will be future pieces of work to pull additional data about the players scores and prices from https://www.guysports.co.uk/guysports/players.php?code=GK but this initial piece of work will focus on the players in the users teams. ideally the data will be stored in yaml files, though if this becomes too large, the ai may suggest a more efficient approach. the useful processed data (initially how many users have selected each player, and how many changes each player had made ) will be presented in such a way that they can be viewed in github and locally. md files with something like mermaid, or github pages seem like the most appropriate. Havnig separate commands/tools for retrieving data and processing and presenting the data seems like a sensible approach."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Capture and retain each user's current team over time (Priority: P1)

A league administrator or analyst needs to understand how team composition changes across a season without relying on live site queries. The system must periodically capture the current team rosters for all Guy Sports users and preserve each captured version as historical data.

**Why this priority**: Historical team data is the foundation for all later analysis, such as player popularity, team changes, and season trends.

**Independent Test**: Can be fully tested by running the periodic capture process for a known set of users and verifying that each captured team snapshot is stored as a distinct historical record.

**Acceptance Scenarios**:

1. **Given** a season page listing active users and their selected players, **When** the scheduled capture runs, **Then** the system records a new snapshot for each user and stores the team composition from that point in time.
2. **Given** a user has changed players since the previous capture, **When** the next capture runs, **Then** the system preserves both the prior and current team states and records the change.
3. **Given** the same user appears on multiple season pages, **When** the capture runs, **Then** the system avoids losing historical records and treats the new data as an additional point in time rather than overwriting the past.

---

### User Story 2 - Show how often each player is selected by users (Priority: P1)

A user wants to see how many managers selected each player across the league. The system must process stored snapshots into an aggregated view that shows player selection counts over time and for the current state.

**Why this priority**: This is the primary analytical output described for the first release and provides useful insight for fans and analysts.

**Independent Test**: Can be fully tested by processing captured team snapshots and verifying that the reported player counts match the number of users holding each player in the selected time window.

**Acceptance Scenarios**:

1. **Given** a set of historical team snapshots, **When** the aggregation step runs, **Then** the system reports how many users selected each player in the most recent snapshot and in the selected comparison window.
2. **Given** a player is selected by multiple users, **When** the report is generated, **Then** the total counts are summed across the stored data without duplicating the same user's selection.

---

### User Story 3 - Show each user's team-change activity across the season (Priority: P1)

A user wants to understand how many team changes each manager has made during the season. The system must track changes across captured snapshots and present each user's change count in a readable summary.

**Why this priority**: This is a core decision-support metric described in the initial feature and directly supports the business value of the project.

**Independent Test**: Can be fully tested by comparing two consecutive user snapshots and verifying that the system counts only actual team differences while preserving the season total for the user.

**Acceptance Scenarios**:

1. **Given** a user has a team snapshot and a later snapshot with different players, **When** the change analysis runs, **Then** the system counts the number of changes made since the previous snapshot.
2. **Given** a user has six or fewer allowed changes in the season, **When** the seasonal summary is generated, **Then** the system reports the total change count in a way that matches the known league rules.
3. **Given** a user has no changes between snapshots, **When** the analysis runs, **Then** the system does not count a change where the roster is unchanged.

---

### User Story 4 - Present results in a GitHub-friendly and local-friendly format (Priority: P2)

An analyst needs a way to review computed results in a lightweight, shareable format without needing a live application. The system must produce human-readable output, such as Markdown with charts or summaries, that can be viewed in GitHub and locally.

**Why this priority**: This enables the project to deliver value even before a full user-facing web experience exists.

**Independent Test**: Can be fully tested by generating a report file and confirming that it renders clearly in a Markdown viewer or GitHub repository page.

**Acceptance Scenarios**:

1. **Given** processed data for player counts and change totals, **When** the presentation step runs, **Then** the system outputs a readable Markdown summary that can be viewed in GitHub and locally.
2. **Given** the output includes chart-like or visual summaries, **When** the document is rendered, **Then** it remains understandable without requiring live data access.

### Edge Cases

- What happens when a user is absent from a season page during a scheduled capture?
- What happens if a user’s team is unchanged between two captures, but a valid historical record must still exist?
- What happens when an existing stored snapshot is incomplete or malformed?
- How does the system handle team pages that include more users than the previously known total?
- How does the system handle duplicate records from multiple capture runs for the same user and time window?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST support scheduled periodic capture of Guy Sports user team data from the season table pages.
- **FR-002**: The system MUST capture the current team roster for each user appearing on the Guy Sports season pages.
- **FR-003**: The system MUST preserve each captured team snapshot as historical data rather than overwriting prior state.
- **FR-004**: The system MUST store team snapshots additively so that past versions remain available for analysis and comparison over time.
- **FR-005**: The system MUST avoid live queries to the original Guy Sports data source during results presentation or user-facing queries.
- **FR-006**: The system MUST treat captured data as the source of truth for user-facing reporting and historical analysis.
- **FR-007**: The system MUST record the players in each user's team for each capture event, including enough metadata to identify the team state at that point in time.
- **FR-008**: The system MUST track changes between successive team snapshots so that user-level change counts can be calculated accurately.
- **FR-009**: The system MUST compute the number of users who selected each player across the relevant captured snapshots and time window.
- **FR-010**: The system MUST compute each user's total team changes within a season using historical team snapshots.
- **FR-011**: The system MUST support processing of captured data into viewable summary outputs for local and GitHub-based review.
- **FR-012**: The system MUST provide a format for presenting processed results in Markdown and visual summaries that do not require live data access.
- **FR-013**: The system MUST separate data capture, data processing, and presentation concerns into distinct commands or tools.
- **FR-014**: The system MUST support manual execution of capture commands for initial setup and testing before automation is enabled.
- **FR-015**: The system MUST explicitly exclude player score data and price data from this initial phase and focus only on players listed in users' teams.
- **FR-016**: The system MUST preserve historical records even when a team remains unchanged between two capture runs, so the record of time and state remains auditable.
- **FR-017**: The system MUST be designed to support future expansion for additional data sources and additional player-related datasets without breaking the historical capture model.

### Key Entities

- **User Team Snapshot**: A record of one user's team state at a specific capture time, including the season context and the roster of selected players.
- **Player Selection**: A player that appears in a user's team snapshot, along with the user and time context that identifies the selection event.
- **Season**: The league season or time window within which team snapshots are compared and aggregated.
- **User Change Count**: A derived value showing how many times a user has changed their team between snapshots during the season.
- **Player Selection Aggregate**: A processed summary showing how many users selected each player across the stored dataset.
- **Presentation Output**: A generated Markdown or chart-based document used to view processed results in GitHub or locally.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The system captures team snapshots for all known users on Guy Sports season pages during each scheduled run, with no silent omission of page entries.
- **SC-002**: Historical records remain available for all processed snapshots, allowing analysts to compare team states across the season without losing prior data.
- **SC-003**: Player selection counts are reported accurately for the current state and selected historical windows, with no double-counting of a user's team membership.
- **SC-004**: Each user’s total team-change count is calculated accurately from sequential snapshots without overcounting unchanged states.
- **SC-005**: Processed outputs are generated in a repo-friendly format that can be reviewed in GitHub and locally without a live data dependency.
- **SC-006**: The initial release supports Guy Sports as the primary source while leaving room for future expansion to additional data sources and derived metrics.

## Assumptions

- The first release focuses on Guy Sports only; DreamTeamFC capture is planned as a future or parallel feature and not required for this v1 scope.
- The system may store data in structured, versioned files such as YAML if the dataset size remains manageable; a more efficient format can be adopted if scale requires it.
- The season table pages are the authoritative source for current user roster data during the initial implementation phase.
- Historical capture can be performed in scheduled runs or manual execution during testing, with the same processing path used in both modes.
- User queries and presentation views rely only on persisted historical data and not on direct access to the upstream website.
- The data collection process is expected to preserve both state and change history, even when users do not change their teams between consecutive captures.
