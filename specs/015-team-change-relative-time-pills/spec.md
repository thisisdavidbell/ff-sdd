---
type: Feature Specification
title: Team Change Relative-Time Pills
description: Simplify Team changes latest-change cells to rounded relative time with interval-specific pills.
tags: [report, rendering, team-changes, relative-time]
status: stable
feature: 015-team-change-relative-time-pills
sdd_approach: streamlined
input_summary: Remove the exact timestamp from Team changes latest-change cells and classify rounded relative times as Hours, Days, Weeks, or Months.
generated: { by: GitHub Copilot, at: 2026-09-06T00:00:00Z }
---

# Requested Behavior

Each populated Latest change cell in the generated Team changes table MUST display only a rounded elapsed-time value and one category pill. It MUST NOT display the calendar date, clock time, or timezone in the collapsed table row.

The category pill MUST use exactly one of these labels, derived from the elapsed time between the latest change timestamp and the report render time:

- `Hours` for elapsed times less than 24 hours.
- `Days` for elapsed times from 24 hours up to, but not including, 7 days.
- `Weeks` for elapsed times from 7 days up to, but not including, 30 days.
- `Months` for elapsed times of 30 days or more.

The rounded elapsed-time value MUST be expressed in the unit corresponding to its pill. The exact date and time of each change MUST remain available in the existing expandable manager-detail view.

## Scope Boundaries

- Change only the generated Team changes Latest change cell and its supporting renderer presentation logic.
- Preserve Team changes sorting, preview controls, expandable manager-detail behavior, event timestamps in expanded details, responsive layout, and the Player ownership Last change column.
- Do not alter capture, processing, YAML schemas, stored timestamps, timezone conversion, or report navigation.

## Affected Artifacts

- `internal/render/html.go`
- `tests/unit/test_render_output_test.go`
- `docs/index.html`
- `ARCHITECTURE.md`
- `specs/015-team-change-relative-time-pills/spec.md`

## Acceptance Checks

- A Team changes Latest change cell contains no exact date, time, or timezone text.
- A latest change less than 24 hours old renders a rounded hour value and an `Hours` pill.
- A latest change from 24 hours through less than 7 days old renders a rounded day value and a `Days` pill.
- A latest change from 7 days through less than 30 days old renders a rounded week value and a `Weeks` pill.
- A latest change at least 30 days old renders a rounded month value and a `Months` pill.
- The `Hours` pill is visually distinct from the `Days` pill and remains the strongest attention cue among the four interval pills.
- The `Hours` pill MUST use fuchsia (`#c026d3`) as its permanent visually distinctive, high-attention treatment.
- Expanded manager-change details still show the existing exact London-formatted timestamp for every event.
- The Player ownership Last change column retains its existing presentation.
- Focused renderer tests pass, the checked-in `docs/index.html` matches renderer output, and `ARCHITECTURE.md` describes the delivered presentation.

## Assumptions And Decisions

- This is a focused generated-report display change with one validation path, so Streamlined SDD applies.
- A month is a fixed 30-day elapsed-time interval for category and rounding purposes; calendar-month arithmetic is out of scope.
- Rounded values use the nearest whole value in their displayed unit. Category selection uses the unrounded elapsed time and the stated half-open thresholds.
- The Hours pill uses a visually distinctive high-attention treatment that does not resemble the blue Days pill; exact color values remain an implementation decision.
- Crimson is rejected because red conventionally signals an error or negative state. The approved fuchsia selection replaces the temporary color prototype control, which MUST be removed from the report.
- Empty or invalid latest-change timestamps retain the renderer's existing fallback behavior rather than receiving a category pill.
- This specification remains draft until explicit user approval. No implementation, generated output, architecture documentation, or unrelated SDD artifact changes are authorized yet.