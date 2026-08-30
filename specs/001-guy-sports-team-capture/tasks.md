# Tasks: Guy Sports Team Capture and Usage Analysis

**Input**: Design documents from `/specs/001-guy-sports-team-capture/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-commands.md

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the Go project skeleton and basic tooling needed for the capture, process, and render pipeline.

- [X] T001 Create the repository structure for Go CLI commands and internal packages under `cmd/`, `internal/`, `data/`, `docs/`, and `tests/`
- [X] T002 Initialize the Go module and project metadata for the feature in the repository root
- [X] T003 [P] Configure Go formatting and validation workflow (`gofmt`, `go test`, and a simple local validation command)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the data model, storage contracts, and shared infrastructure before user-story work can begin.

**Checkpoint**: Foundation ready — capture, processing, and render work can begin in parallel.

- [X] T004 Define the shared data model for snapshots, player references, and seasonal storage in `internal/models/`
- [X] T005 [P] Implement the storage abstraction for append-only YAML records and season-based directories in `internal/storage/`
- [X] T006 [P] Implement common validation and timestamp helpers used by capture, processing, and rendering in `internal/validation/`
- [X] T007 Create the command scaffolding for `cmd/capture`, `cmd/process`, and `cmd/render` with shared CLI wiring
- [X] T008 Create the seasonal data-layout contract and path resolution helpers for `data/<season>/raw`, `data/<season>/processed`, and `docs/`

---

## Phase 3: User Story 1 - Capture and retain each manager's current team over time (Priority: P1) 🎯 MVP

**Goal**: Capture and preserve timestamped Guy Sports team snapshots in append-only historical storage.

**Independent Test**: Run the capture command against fixture data and verify that each manager snapshot is saved as a distinct historical record without overwriting prior data.

### Tests for User Story 1

- [X] T009 [P] [US1] Add a failing unit test for snapshot serialization and validation in `tests/unit/test_snapshot_model.go`
- [X] T010 [P] [US1] Add a failing integration test for the append-only capture flow in `tests/integration/test_capture_flow.go`

### Implementation for User Story 1

- [X] T011 [P] [US1] Implement `ManagerSnapshot` and `PlayerReference` structs with `team_name`, `manager_name`, and underscore space replacement in `internal/models/snapshot.go`
- [X] T012 [P] [US1] Implement the raw YAML writer/reader for per-manager directories `data/<season>/raw/<team_name>_<manager_id>/<timestamp>.yaml` in `internal/storage/yaml_store.go`
- [X] T013 [US1] Implement multi-page season-table parsing across configured pages (default 3) and manager roster extraction in `internal/capture/parser.go`
- [X] T014 [US1] Implement the capture pipeline that collects manager teams, attaches timestamps, skips duplicate snapshots for unchanged teams, and writes raw YAML under `data/<season>/raw/<team_name>_<manager_id>/`
- [X] T015 [US1] Implement `cmd/capture/main.go` loading configuration from `config.yaml` with the CLI contract for `./bin/capture --source guysports`
- [X] T016 [US1] Add validation for malformed snapshots, partial page failures, and duplicate capture handling

**Checkpoint**: At this point, User Story 1 should be fully functional and independently testable.

---

## Phase 4: User Story 2 - Show how often each player is selected by managers (Priority: P1)

**Goal**: Produce two complementary outputs from the stored raw snapshots into a single human-readable file `data/<season>/processed/player-ownership.yaml`:
1. Use each manager's latest recorded team state to calculate the current player-usage count, counting each player only from the newest team snapshot for that manager.
2. Preserve the per-capture historical player counts across the season so the trend in ownership over time can be graphed and compared.

**Independent Test**: Process known fixture snapshots and confirm that the latest team snapshot for each manager is used to produce the current counts, while the full historical capture sequence is retained to produce the ownership trend over time in `player-ownership.yaml`.

### Tests for User Story 2

- [X] T017 [P] [US2] Add a failing unit test for aggregation logic in `tests/unit/test_player_ownership.go`
- [X] T018 [P] [US2] Add a failing integration test for processing raw snapshots into derived ownership YAML in `tests/integration/test_process_ownership.go`

### Implementation for User Story 2

- [X] T019 [P] [US2] Define the derived aggregate model for player ownership and historical trend in `internal/models/ownership.go`
- [X] T020 [US2] Implement the processing step that reads raw snapshots across manager directories, builds the latest ownership view and historical trend, and writes single output `data/<season>/processed/player-ownership.yaml`
- [X] T021 [US2] Implement ordering and trend logic separately for the current snapshot-based counts and the historical time-series counts
- [X] T021b [US2] Fix point-in-time player ownership logic in `internal/processing/ownership.go` so that counts at timestamp T evaluate active snapshots for ALL managers in the league as of T, preserving counts during partial capture runs where unchanged teams do not produce new snapshot files
- [X] T021c [US2] Add unit test in `tests/unit/test_player_ownership_test.go` verifying that a partial capture run (where only 2 out of 28 managers change) preserves unchanged player counts across current and historical trends
- [X] T022 [US2] Implement `cmd/process/main.go` to run the offline processing workflow from persisted raw data
- [X] T023 [US2] Add clear error handling for missing data, malformed records, and empty seasons

**Checkpoint**: At this point, player selection output should be generated from historical snapshots and be ready for presentation.

---

## Phase 5: User Story 3 - Show each manager's team-change activity across the season (Priority: P1)

**Goal**: Compute manager change counts and change-event details from sequential team snapshots and save them by team name.

**Independent Test**: Compare a manager's consecutive snapshots and verify the total change count, latest change timestamp, and event detail match the expected diffs.

### Tests for User Story 3

- [X] T024 [P] [US3] Add a failing unit test for sequential diffing of manager snapshots in `tests/unit/test_manager_changes.go`
- [X] T025 [P] [US3] Add a failing integration test for manager summary generation in `tests/integration/test_manager_changes.go`

### Implementation for User Story 3

- [X] T026 [P] [US3] Define the derived change-event and manager summary models with `manager_name` and `team_name` in `internal/models/change_summary.go`
- [X] T027 [US3] Implement the diff logic that compares consecutive snapshots and records added/removed players without overcounting unchanged states
- [X] T028 [US3] Implement cumulative season summaries for `total_changes`, `latest_change_at`, and `changed_since_last_snapshot`
- [X] T028b [US3] Fix transfer change count logic in `internal/processing/ownership.go` so single substitutions (1 player added, 1 removed) count as 1 change (`max(added, removed)`), include player names on change events, and calculate `total_changes` as the sum of event change counts
- [X] T028c [US3] Add unit test in `tests/unit/test_manager_changes_test.go` verifying that swapping 1 player out for 1 player in produces `ChangeCount == 1` and `TotalChanges == 1` with player names populated
- [X] T029 [US3] Write derived manager-change YAML output named `<team_name>_<manager_id>.yaml` under `data/<season>/processed/manager-changes/`
- [X] T030 [US3] Add validation to ensure duplicate snapshots and unchanged states are treated correctly

**Checkpoint**: At this point, team-change history should be fully derived and ready for viewing in the static report.

---

## Phase 6: User Story 4 - Present results in a GitHub-friendly and local-friendly format (Priority: P2)

**Goal**: Generate a static HTML report from the latest processed data currently available in the local repository state. The HTML must display player ownership sorted highest count first, manager changes with Manager Name and Team Name, and historical trends as a large line chart with consistent x-axis date intervals and ample vertical height for multi-player readability.

**Independent Test**: Render the report from fixture processed data and verify the HTML is readable locally, contains sorted ownership and change summaries, and renders a large historical trends line chart with equidistant x-axis date spacing and generous height.

### Tests for User Story 4

- [X] T031 [P] [US4] Add a failing unit test for HTML template output, descending ordering, manager/team name display, and historical trends line chart generation in `tests/unit/test_render_output.go`
- [X] T032 [P] [US4] Add a failing integration test for render-from-processed-data flow in `tests/integration/test_render_flow.go`

### Implementation for User Story 4

- [X] T033 [P] [US4] Implement the HTML render model, line chart generation with consistent x-axis spacing and comfortable vertical scale, and template structure in `internal/render/`
- [X] T034 [US4] Implement the render pipeline that reads processed `player-ownership.yaml` and manager change YAML files and emits static HTML with the large historical trends line chart into `docs/index.html`
- [X] T035 [US4] Implement `cmd/render/main.go` to generate the latest report with descending player ownership sort, Manager/Team name display, and the large historical trends line chart
- [X] T035b [US4] Implement interactive collapsible drill-down for Manager Changes in `internal/render/html.go` displaying event timestamps, `+ Added Player Name`, `- Removed Player Name`, and event change counts
- [X] T035c [US4] Add unit test in `tests/unit/test_render_output_test.go` verifying interactive toggle script, detail row markup, and added/removed player names in HTML render
- [X] T036 [US4] Ensure the rendered report never queries live sources and preserves the last known-good HTML if local generation fails
- [X] T037 [US4] Add a reset/local cleanup guide and validation notes to the quickstart workflow, including the rule that uncommitted local processed/render changes can be discarded instead of being committed
- [X] T037b [P] Document the optional Guy Sports live smoke test as a narrowly scoped check for retrieval-code changes or explicit requests only; it is not part of the normal validation cycle

**Checkpoint**: The system should now be able to show static, reader-friendly results from historical stored data without network access.

---

## Phase 7: Cross-Cutting Validation and Polish

**Purpose**: Final repository-wide verification, documentation, and offline smoke checks.

- [X] T038 [P] Run the offline Go test suite across unit and integration tests for capture, processing, and rendering
- [X] T039 [P] Validate the CLI commands against fixture datasets and confirm output paths and file formats match the planned storage model
- [X] T040 [P] Update the project README and feature quickstart to explain the capture → process → render workflow and reset policy
- [X] T041 [P] Verify the repository honors the governance constraints: append-only raw data, committed processed-data authority, and no live-source validation by default
- [X] T042 Run the full offline validation command from the quickstart and confirm the generated HTML is consistent with the committed processed state

---

## Dependencies & Execution Order

### Phase Dependencies

- Phase 1 Setup: no dependencies
- Phase 2 Foundational: depends on Setup completion and blocks all user story work
- Phase 3+ User Stories: all depend on Phase 2 completion
- Phase 7 Polish: depends on all user stories being complete

### User Story Dependencies

- User Story 1 (P1): capture data archive first and must complete before any processing or render work can rely on raw data
- User Story 2 (P1): depends on User Story 1 data being available
- User Story 3 (P1): depends on User Story 1 raw snapshots and change comparison logic
- User Story 4 (P2): depends on User Story 2 and User Story 3 processed data

### Parallel Opportunities

- Setup tasks and foundational tasks marked [P] can run in parallel
- Tests for each user story can be implemented in parallel with their corresponding story logic once the foundation is ready
- Different stories may proceed in parallel if the team has capacity, but Story 2 and Story 3 rely on Story 1 captured data paths and storage conventions

---

## Notes

- [P] tasks are independent and can be executed in parallel where files do not overlap
- All routine validation must remain offline unless a narrow live smoke test is explicitly approved for Guy Sports retrieval code changes or a direct request
- The append-only raw archive is the source of truth and must not be overwritten under normal operation
- The latest processed output currently present in the local repository is the input for the current HTML render; it may be the latest committed state or the latest state from a local capture/process run, and the render step does not need to know which it is
- Uncommitted local processed/render changes may be discarded without being committed, while the checked-in source data remains authoritative
- The implementation intentionally keeps the workflow simple, file-based, and human-browsable for the initial release
