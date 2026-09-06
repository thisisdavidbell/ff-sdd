---
type: Feature Specification
title: "Render Change Context"
description: Historical spec artifact for 007-render-change-context.
tags: [sdd, feature-specification, 007-render-change-context]
status: stable
feature: 007-render-change-context
sdd_approach: streamlined
input_summary: Historical spec artifact.
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
---

# Feature Specification: Render Change Context


## User Scenarios & Testing

### User Story 1 - Prioritize recent manager activity (Priority: P1)

As a report reader, I want manager changes ordered from the most recently updated manager to the least recently updated so that I can see the newest activity immediately.

**Acceptance Scenarios**:

1. **Given** manager summaries with different latest-change times, **When** the report is rendered, **Then** the newest change is listed first.
2. **Given** manager summaries with no changes, **When** the report is rendered, **Then** they appear below all managers with a latest-change time.

### User Story 2 - Understand ownership movement (Priority: P1)

As a report reader, I want each ownership count to indicate its latest direction and the time it last changed so that I can interpret the current total in context.

**Acceptance Scenarios**:

1. **Given** a player's latest ownership change increased its count, **When** the report is rendered, **Then** the count has a green up arrow and the last-change time is shown.
2. **Given** a player's latest ownership change decreased its count, **When** the report is rendered, **Then** the count has a red down arrow and the last-change time is shown.
3. **Given** a player has no ownership-count change, **When** the report is rendered, **Then** no direction icon or last-change value is shown.

### User Story 3 - Orient new readers (Priority: P2)

As a human reader, I want a concise README overview and links to detailed documentation so that I can understand the repository and find deeper technical context quickly.

**Acceptance Scenarios**:

1. **Given** a reader opens the README, **When** they scan it, **Then** they can identify the repository's purpose, pipeline stages, report output, and links to deeper documentation.

### Edge Cases

- Invalid timestamps retain a deterministic sort order and do not prevent report rendering.
- A player history containing repeated equal counts does not display a direction until a count change exists.

## Requirements

### Functional Requirements

- **FR-001**: The Manager changes table MUST order summaries with valid latest-change times newest first and summaries without a latest-change time last.
- **FR-002**: The Player ownership table MUST include a Last change column.
- **FR-003**: The ownership count MUST display a green up arrow after its most recent increase and a red down arrow after its most recent decrease.
- **FR-004**: The Last change column and direction indicator MUST be blank for a player with no ownership-count change.
- **FR-005**: The README MUST concisely explain the application's goals, capture-process-render capability, generated report, main commands, and deeper documentation links.
- **FR-006**: `ARCHITECTURE.md` MUST describe the report's ownership-change context and latest-activity ordering.

## Success Criteria

- **SC-001**: Rendering mixed manager summaries places every timestamped latest change above summaries with no latest change.
- **SC-002**: Rendering ownership histories accurately shows the latest increase, decrease, or no-change state for every player.
- **SC-003**: A reader can identify the application's purpose, pipeline, output, and detailed documentation from the README without reading source code.

## Assumptions

- Ownership history entries represent chronological point-in-time counts.
- Existing timestamp formatting remains appropriate for the new last-change column.
- The change is limited to the generated report and repository documentation; no input-data schema changes are needed.