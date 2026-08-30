# Tasks: Automated Data Pipeline & Capture Source Separation

**Input**: Design documents from `/specs/002-automated-data-pipeline/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-commands.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story?] Description with file path`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- File paths are included in all task descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify development environment and setup workspace structure

- [X] T001 Verify project structure and Go dependencies per `specs/002-automated-data-pipeline/plan.md`

---

## Phase 2: User Story 3 - Source-Specific Data Capture Command Architecture (Priority: P2)

**Goal**: Structure the capture entrypoint so GuySports collection is explicitly isolated in `cmd/capture-guysports`, preparing for future alternative sources (such as DreamTeam).

**Independent Test**: Execute `go run ./cmd/capture-guysports` directly and verify raw YAML snapshots are fetched and written to `data/2026-27/raw/`.

### Implementation for User Story 3

- [X] T002 [US3] Move `cmd/capture/` directory to `cmd/capture-guysports/`
- [X] T003 [P] [US3] Update capture target references and test imports in `tests/integration/test_capture_flow_test.go` and `tests/unit/test_capture_parser_test.go`
- [X] T004 [P] [US3] Update `README.md` to document `go run ./cmd/capture-guysports` usage

**Checkpoint**: User Story 3 complete - `cmd/capture-guysports` is isolated and verified independently.

---

## Phase 3: User Story 1 - Full End-to-End Pipeline Automation (Priority: P1) 🎯 MVP

**Goal**: Create root pipeline execution script (`run.sh`) that invokes `go run ./cmd/capture-guysports`, `go run ./cmd/process`, and `go run ./cmd/render` sequentially with strict error handling (`set -euo pipefail`).

**Independent Test**: Execute `./run.sh` locally and verify sequential execution through capture, process, and render via `go run`, halting immediately if any stage fails.

### Tests for User Story 1

- [X] T005 [P] [US1] Create script verification test in `tests/unit/test_run_script_test.go` to validate error trapping and stage ordering of `run.sh`

### Implementation for User Story 1

- [X] T006 [US1] Implement root pipeline execution script `run.sh` with `set -euo pipefail` sequentially invoking `go run ./cmd/capture-guysports`, `go run ./cmd/process`, and `go run ./cmd/render`
- [X] T007 [US1] Grant executable permissions to `run.sh` (`chmod +x run.sh`) and add inline usage comments

**Checkpoint**: User Story 1 complete - `./run.sh` provides one-command local execution using `go run`.

---

## Phase 4: User Story 2 - Scheduled & On-Demand CI/CD Automation (Priority: P1)

**Goal**: Configure GitHub Actions workflow to execute `./run.sh` weekly on Fridays at 18:00 UTC and on manual dispatch, committing updated data artifacts back to the repository branch.

**Independent Test**: Validate GitHub Actions workflow definition syntax for schedule (`0 18 * * 5`), manual trigger (`workflow_dispatch`), script execution (`./run.sh`), and git commit/push steps.

### Implementation for User Story 2

- [X] T008 [US2] Create GitHub Actions workflow file `.github/workflows/schedule-run-action.yml` with `schedule` (`cron: '0 18 * * 5'`) and `workflow_dispatch` triggers
- [X] T009 [US2] Add workflow build, checkout, Go setup, `./run.sh` execution, and git commit (`Automated weekly pipeline run`) & push steps in `.github/workflows/schedule-run-action.yml`

**Checkpoint**: User Story 2 complete - scheduled and manual CI/CD automation configured.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final documentation updates and manual verification guidance

- [X] T010 [P] Update `AGENTS.md` repository guidelines to reflect `cmd/capture-guysports` entrypoint and `run.sh` automation
- [X] T011 Prompt user to manually run and verify quickstart validation scenarios in `specs/002-automated-data-pipeline/quickstart.md` and clean up working tree changes afterwards

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **User Story 3 (Phase 2)**: Can start after Setup — renames `cmd/capture` to `cmd/capture-guysports`
- **User Story 1 (Phase 3)**: Depends on User Story 3 (`cmd/capture-guysports` location)
- **User Story 2 (Phase 4)**: Depends on User Story 1 (`run.sh` script)
- **Polish (Phase 5)**: Depends on all user stories being complete

---

## Parallel Opportunities

- T003 (unit/integration test updates) and T004 (`README.md` updates) can run in parallel after T002 directory move
- T005 (script verification test) can run in parallel with T006 (`run.sh` creation)
- T010 (`AGENTS.md` updates) can run in parallel with polish tasks

---

## Implementation Strategy

### MVP First

1. Complete Phase 1: Setup
2. Complete Phase 2: User Story 3 (directory rename to `cmd/capture-guysports`)
3. Complete Phase 3: User Story 1 (`run.sh` using `go run`)
4. **STOP and VALIDATE**: Test `run.sh` locally
5. Complete Phase 4: User Story 2 (GitHub Actions workflow)
