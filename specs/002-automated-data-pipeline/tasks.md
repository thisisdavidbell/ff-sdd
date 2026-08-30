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

- [ ] T001 Verify project structure and Go dependencies per `specs/002-automated-data-pipeline/plan.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core source restructuring required before pipeline execution scripts can target `capture-guysports`

**⚠️ CRITICAL**: Must complete before User Story pipeline scripts call `capture-guysports`

- [ ] T002 Move `cmd/capture/` directory to `cmd/capture-guysports/` in `cmd/capture-guysports/main.go`
- [ ] T003 [P] Update capture target references and test paths in `tests/integration/test_capture_flow_test.go` and `tests/unit/test_capture_parser_test.go`

**Checkpoint**: Foundation ready - `capture-guysports` source target is isolated and testable

---

## Phase 3: User Story 1 - Full End-to-End Pipeline Automation (Priority: P1) 🎯 MVP

**Goal**: Create root pipeline execution script (`run.sh`) that builds binaries and runs capture, process, and render stages sequentially with strict error handling (`set -euo pipefail`).

**Independent Test**: Execute `run.sh` locally (or in test mode) and verify sequential execution through capture, process, and render, stopping immediately if any stage fails.

### Tests for User Story 1

- [ ] T004 [P] [US1] Create script verification test in `tests/unit/test_run_script_test.go` to validate error trapping and stage ordering of `run.sh`

### Implementation for User Story 1

- [ ] T005 [US1] Implement root pipeline execution script `run.sh` with `set -euo pipefail` sequentially invoking `capture-guysports`, `process`, and `render`
- [ ] T006 [US1] Grant executable permissions to `run.sh` (`chmod +x run.sh`) and add inline usage documentation

**Checkpoint**: At this point, User Story 1 is fully functional and can be tested locally via `./run.sh`

---

## Phase 4: User Story 2 - Scheduled & On-Demand CI/CD Automation (Priority: P1)

**Goal**: Configure GitHub Actions workflow to run `./run.sh` automatically every Friday at 18:00 UTC and on manual dispatch, committing updated data artifacts back to the repository branch.

**Independent Test**: Validate GitHub Actions workflow definition syntax for schedule (`0 18 * * 5`), manual trigger (`workflow_dispatch`), script execution (`./run.sh`), and git commit/push steps.

### Implementation for User Story 2

- [ ] T007 [US2] Create GitHub Actions workflow file `.github/workflows/schedule-run-action.yml` with `schedule` (`cron: '0 18 * * 5'`) and `workflow_dispatch` triggers
- [ ] T008 [US2] Add workflow build, checkout, Go setup, `./run.sh` execution, and git commit (`Automated weekly pipeline run`) & push steps in `.github/workflows/schedule-run-action.yml`

**Checkpoint**: User Stories 1 AND 2 are complete - local execution and scheduled CI/CD workflow are fully configured

---

## Phase 5: User Story 3 - Source-Specific Data Capture Command Architecture (Priority: P2)

**Goal**: Refine `capture-guysports` CLI interface and project documentation to isolate GuySports capture from future data sources.

**Independent Test**: Build and run `bin/capture-guysports` directly and verify flag parsing (`-config`, `-pages`, `-source`) operates independently.

### Implementation for User Story 3

- [ ] T009 [P] [US3] Update CLI flag defaults and help text in `cmd/capture-guysports/main.go` to explicitly identify GuySports provider isolation
- [ ] T010 [P] [US3] Update `README.md` to document `capture-guysports` CLI usage and `run.sh` automated execution workflow

**Checkpoint**: All user stories are independently functional and documented

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final documentation updates and manual verification guidance

- [ ] T011 [P] Update `AGENTS.md` repository guidelines to reflect `capture-guysports` entrypoint and `run.sh` automation
- [ ] T012 Prompt user to manually run and verify quickstart validation scenarios in `specs/002-automated-data-pipeline/quickstart.md` and clean up working tree changes afterwards

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User Story 1 (P1) -> User Story 2 (P1) -> User Story 3 (P2)
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Foundational (Phase 2) for `cmd/capture-guysports` target
- **User Story 2 (P1)**: Depends on User Story 1 (`run.sh` script)
- **User Story 3 (P2)**: Depends on Foundational (Phase 2); can run in parallel with US2

---

## Parallel Opportunities

- T003 (unit/integration test updates) can run in parallel with T002 once directory is moved
- T004 (script verification test) can run in parallel with T005
- T009 (`cmd/capture-guysports/main.go` flag updates) and T010 (`README.md` updates) can run in parallel
- T011 (`AGENTS.md` updates) can run in parallel with polish tasks

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (`cmd/capture` -> `cmd/capture-guysports`)
3. Complete Phase 3: User Story 1 (`run.sh`)
4. **STOP and VALIDATE**: Test `run.sh` locally
5. Continue to User Story 2 (GitHub Actions) and User Story 3 (Documentation & Refinements)
