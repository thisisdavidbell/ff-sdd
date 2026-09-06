---
type: Feature Specification
title: "London Timestamps"
description: Historical spec artifact for 009-london-timestamps.
tags: [sdd, feature-specification, 009-london-timestamps]
status: stable
feature: 009-london-timestamps
sdd_approach: streamlined
input_summary: Historical spec artifact.
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
---

# Feature Specification: London Timestamps


## User Scenarios & Testing

### User Story 1 - Read unambiguous report times (Priority: P1)

As a report reader, I want every displayed report time in UK local time and minute precision so that I can interpret activity without converting a timezone or parsing technical timestamp notation.

**Why this priority**: The report's value depends on readers accurately understanding when data and changes occurred.

**Independent Test**: Render UTC report data from dates in both UK winter and summer time and verify every visible date and time is London-local, omits seconds, and contains no `Z` suffix.

**Acceptance Scenarios**:

1. **Given** a UTC source timestamp from a date in British Summer Time, **When** the report displays it, **Then** the displayed time is one hour ahead of UTC and contains only the date and hour-and-minute.
2. **Given** a UTC source timestamp from a date in Greenwich Mean Time, **When** the report displays it, **Then** the displayed time matches UK winter time and contains only the date and hour-and-minute.
3. **Given** a report is generated, **When** a reader views its generation time, **Then** it is shown in London time without seconds or a `Z` suffix.

### User Story 2 - Receive scheduled updates at a predictable UK time (Priority: P1)

As a report reader, I want the automated data pipeline to run daily at 01:54 London time so that updates avoid the hourly GitHub Actions peak and retain their local-time schedule through daylight-saving changes.

**Why this priority**: Reliable scheduled collection preserves the historical data record and avoids a known period of elevated scheduling load.

**Independent Test**: Inspect the workflow schedule and verify it specifies the `Europe/London` timezone and a daily 01:54 run time.

**Acceptance Scenarios**:

1. **Given** the daily schedule is evaluated in winter, **When** 01:54 in London occurs, **Then** the pipeline is scheduled to run.
2. **Given** the daily schedule is evaluated in summer, **When** 01:54 in London occurs, **Then** the pipeline is scheduled to run.

### User Story 3 - Preserve a stable timestamp record (Priority: P2)

As a maintainer, I want captured and processed timestamp records to remain in UTC so that historical data has one stable timezone representation and can support future displays in other timezones.

**Why this priority**: UTC avoids seasonal changes in stored data while preserving an unambiguous, portable historical record.

**Independent Test**: Run or unit-test timestamp creation in both UK winter and summer and verify each persisted value is a valid UTC timestamp.

**Acceptance Scenarios**:

1. **Given** a capture runs in British Summer Time, **When** it records its capture time, **Then** the stored timestamp remains in UTC.
2. **Given** processing runs in Greenwich Mean Time, **When** it records its generation time, **Then** the stored timestamp remains in UTC.

### Edge Cases

- UTC timestamps remain the canonical stored format and display as their equivalent London local time.
- A timestamp during the UK daylight-saving transition is converted using the London timezone rules rather than a fixed offset.
- Invalid stored timestamps retain the renderer's current graceful fallback behavior and do not stop report generation.
- The renderer fails rather than publishing a report with an incorrect fallback timezone when London timezone data is unavailable.
- The 01:54 schedule is before the UK spring-forward skipped hour, so it remains a valid daily local time.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST create all persisted capture and processing timestamps in UTC.
- **FR-002**: The report MUST convert every valid displayed timestamp, including report generation, ownership changes, manager changes, event history, chart labels, and timestamp-based hover text, to `Europe/London` before display.
- **FR-003**: Every valid timestamp displayed in the report MUST use date-and-minute precision and MUST NOT show seconds or a `Z` suffix.
- **FR-004**: The scheduled pipeline workflow MUST run daily at 01:54 in the `Europe/London` timezone.
- **FR-005**: The scheduled pipeline workflow MUST use its scheduler's native timezone support so daylight-saving transitions retain the intended local schedule without duplicate schedules or a custom time guard.
- **FR-006**: Timestamp parsing and chronological ordering MUST continue to accept valid RFC3339 timestamps, including canonical UTC data and timestamps with fractional seconds.
- **FR-007**: Invalid timestamps MUST continue to avoid preventing report rendering.
- **FR-008**: `ARCHITECTURE.md` MUST describe the London-time timestamp and scheduling behavior.
- **FR-009**: The renderer MUST fail rather than publish report timestamps in another timezone when `Europe/London` timezone data is unavailable.

### Key Entities

- **Capture timestamp**: The UTC time at which a manager snapshot is collected.
- **Processing timestamp**: The UTC time at which ownership output is produced.
- **Report timestamp**: A reader-facing minute-precision London-time conversion of a stored UTC timestamp.
- **Scheduled pipeline run**: The daily automated collection and report-generation event, defined in London local time.

## Success Criteria

- **SC-001**: All valid timestamps shown in a rendered report are London-local and match the `YYYY-MM-DD HH:MM` shape, with zero occurrences of seconds or a `Z` suffix in displayed timestamps.
- **SC-002**: Test data spanning UK winter and summer produces correct London local times for 100% of visible timestamp locations.
- **SC-003**: The workflow configuration schedules exactly one daily pipeline run at 01:54 London time in both Greenwich Mean Time and British Summer Time.
- **SC-004**: All persisted timestamps remain in UTC while valid UTC timestamp records render without failure.

## Assumptions

- UTC is the authoritative persisted timestamp format; `Europe/London` and its daylight-saving rules apply only when rendering reader-facing times.
- The requested off-hour time is 01:54 London time every day.
- Machine-readable persisted timestamps retain RFC3339 precision in UTC; the removal of seconds applies only to reader-facing HTML.
- The report's existing behavior for invalid timestamps remains the expected fallback behavior.
- The render environment provides IANA timezone data for `Europe/London`; rendering fails clearly when that dependency is unavailable.
- GitHub Actions native `timezone` support is available for this repository, as documented by GitHub on 2026-09-05.