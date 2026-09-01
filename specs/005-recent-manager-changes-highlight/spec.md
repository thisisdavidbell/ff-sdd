# Feature Specification: Recent Manager Change Highlight

**Created**: 2026-09-01
**Status**: Approach decided — ready for implementation
**Input**: User description: "Update the rendering of the Manager Changes table to give a visual indication of changes that have happened recently, likely next to the Latest change date."

## User Scenarios & Testing

### User Story 1 - Spot Recent Activity at a Glance (Priority: P1)

As a report viewer, I want to immediately see which managers have made changes recently so that I can focus on the most current activity without reading every "Latest change" date.

**Why this priority**: The Manager Changes table already lists a latest-change date per manager, but a viewer must read and mentally compare each date to "now" to judge recency. A visual cue removes that effort and is the entire point of this feature.

**Independent Test**: Generate a report where at least one manager's latest change falls within each of the day/week/month tiers and at least one falls outside all tiers. Confirm each row shows the correct relative-time phrase and, where applicable, the correct single tiered pill, without needing to read the raw date.

**Acceptance Scenarios**:

1. **Given** a manager's latest change occurred, **When** a viewer looks at the Manager Changes table, **Then** the "Latest change" cell shows the absolute date plus a relative-time phrase (e.g., "2 hours ago").
2. **Given** a manager's latest change falls within the last day, week, or month, **When** a viewer looks at the table, **Then** that row additionally shows a tiered pill labeled "Day", "Week", or "Month" (whichever is the most specific matching tier).
3. **Given** a manager's latest change is older than a month, or a manager has no recorded changes, **When** a viewer looks at the table, **Then** no tiered pill is shown (the relative-time phrase, if any, is still shown).
4. **Given** the report is regenerated at a later date, **When** a previously shown tier ages into a broader tier (or out of all tiers), **Then** the pill updates accordingly (or disappears) and the relative-time phrase updates.
5. **Given** the viewer is using dark mode, **When** the relative-time text and tiered pill are shown, **Then** both remain clearly visible and legible in both light and dark themes.

### Edge Cases

- A manager with no recorded changes (`TotalChanges` is 0 / `LatestChangeAt` empty) must not show a relative-time phrase or a tiered pill.
- The tiered pill must not be confused with the existing row-expand affordance (▶) already used for rows with change history.
- Behavior at tier boundaries (e.g., a change exactly at 24 hours, 7 days, or 30 days) must be deterministic and consistently applied; each row shows only its single most-specific tier, never more than one pill.
- The tiered pills rely on color plus a text label ("Day"/"Week"/"Month"), so they remain understandable for viewers who cannot perceive color alone.

## Requirements

### Functional Requirements

- **FR-001**: The Manager Changes table MUST show, within the "Latest change" cell, a relative-time phrase (e.g., "2 hours ago") alongside the existing absolute date, for any manager with a recorded change.
- **FR-002**: The Manager Changes table MUST additionally show a tiered pill labeled "Day", "Week", or "Month" next to the date when a manager's latest change falls within that tier (< 24 hours, < 7 days, < 30 days respectively, most-specific tier wins).
- **FR-003**: Rows whose latest change is older than a month, or managers with no recorded changes, MUST NOT display a tiered pill. Managers with no recorded changes MUST also NOT display a relative-time phrase.
- **FR-004**: The relative-time phrase and tiered pills MUST be distinguishable in both light and dark report themes.
- **FR-005**: The tiered pill MUST NOT rely on color alone to convey meaning; each pill MUST include a text label ("Day"/"Week"/"Month") in addition to its color.
- **FR-006**: This feature MUST NOT alter underlying manager-change data, only its presentation.
- **FR-007**: Each row MUST show at most one tiered pill (the most specific matching tier), never multiple pills simultaneously.

## Candidate Approaches Considered

The following options were prototyped directly in `docs/index.html` for visual comparison:

1. **Badge/pill label** — A small "New" or "Recent" pill next to the date (e.g., `2026-08-30 14:02  🆕`). Catches the eye but doesn't convey how recent.
2. **Icon marker** — A single glyph next to the date (e.g., 🔥, ⭐, ●) with a title/tooltip explaining the recency window.
3. **Row or cell background tint** — Subtle background highlight on the "Latest change" cell (or whole row) for recent entries, paired with a non-color cue.
4. **Bold/emphasized date text** — Render the date itself in bold (or a distinct weight/underline) when recent, with an accessible label (e.g., `aria-label="recent change"`).
5. **Relative time supplement** — Show a relative-time phrase alongside the absolute date (e.g., "2 days ago"). Gives useful detail but doesn't catch the eye on its own.
6. **Tiered/graduated indicator (dots)** — Multiple thresholds with different dot counts (e.g., ●●● within 24h, ●● within 3 days, ● within 7 days) to convey "how recent," not just "recent or not."
7. **Animated pulse/glow** — A subtle CSS animation (e.g., pulsing dot) drawing attention to very recent changes; higher visual novelty, more implementation/testing overhead, and should still degrade gracefully (e.g., `prefers-reduced-motion`).

## Decision

**Chosen approach**: Option 5 (relative-time supplement) combined with a series of tiered pills, labeled "Day", "Week", and "Month", covering changes within the last 24 hours, 7 days, and 30 days respectively. This was chosen because relative-time text alone gives useful detail but doesn't catch the eye, while a plain "new" pill (option 1) catches the eye but loses the "how recent" detail — the combination gives both. Each row shows its relative-time phrase plus, if applicable, the single most-specific tiered pill; no pill is shown once a change is older than a month.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A viewer can identify all managers with recent changes in the Manager Changes table within a few seconds, without reading individual dates, via the tiered pills.
- **SC-002**: 100% of rows with a latest change inside a tier (day/week/month) display exactly the correct single pill; 100% of rows outside all tiers (or with no changes) display no pill.
- **SC-003**: 100% of rows with a recorded change display an accurate relative-time phrase alongside the absolute date.
- **SC-004**: The relative-time phrase and pills remain legible and distinguishable in both light and dark themes, and pills do not depend on color alone.
- **SC-005**: Regenerating the report at a later time correctly updates or removes the pill and relative-time phrase as each change ages through the day/week/month tiers.

## Assumptions

- Recency (relative-time and tier calculations) is evaluated relative to the report's generation time (already rendered in the report header), not the viewer's local clock.
- `LatestChangeAt` (RFC3339 string) is the field used to determine recency; this feature does not need to inspect individual `EventHistory` entries beyond the existing latest-change value.
- "Day" = within 24 hours, "Week" = within 7 days, "Month" = within 30 days, evaluated in that priority order so only the single most-specific tier's pill is shown.

## Validation

- Prototyped in `docs/index.html` (2026-09-01): relative-time text plus tiered Day/Week/Month pills confirmed as the chosen visual approach.
- Remaining: remove the prototype section from `docs/index.html`, implement relative-time formatting and tiered pills in `internal/render/html.go`, then run `go test ./...` and manually verify light/dark themes.
