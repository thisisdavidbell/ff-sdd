---
type: Feature Specification
title: Mobile Team Change Expansion
description: Restore visible expandable Team changes details when a summary row is tapped on a narrow viewport.
tags: [report, mobile, team-changes, rendering, bug-fix]
status: stable
feature: 017-mobile-team-change-expansion
sdd_approach: streamlined
input_summary: Fix Team changes rows that toggle correctly on touch devices but do not visibly reveal their detail content.
generated: { by: GitHub Copilot, at: 2026-09-06T00:00:00Z }
---

# Requested Behavior

On a narrow viewport, tapping a Team changes summary row with event history MUST reveal its existing manager-change detail row. Tapping the same summary row again MUST collapse that detail row.

The mobile Manager-column rule MUST hide only the Team changes header and Manager cells in summary rows. It MUST NOT hide the single spanning cell in a `manager-detail-row`.

## Scope Boundaries

- Change only generated report HTML/CSS/browser behavior and focused renderer tests needed to restore mobile detail visibility.
- Preserve desktop Team changes expansion behavior, table sorting, preview controls, hidden Manager column on mobile, detail content, timestamps, Player ownership presentation, and report navigation.
- Do not change capture, processing, YAML schemas, stored timestamps, external dependencies, or add a server-side endpoint.

## Affected Artifacts

- `internal/render/html.go`
- `tests/unit/test_render_output_test.go`
- `docs/index.html`
- `ARCHITECTURE.md`
- `specs/017-mobile-team-change-expansion/spec.md`
- `specs/index.md`
- `specs/log.md`

## Acceptance Checks

- At a narrow viewport, the Manager header and Manager cell of each Team changes summary row are hidden.
- At a narrow viewport, tapping a Team changes summary row with event history visibly reveals its existing detail content.
- At a narrow viewport, tapping the same summary row again collapses its detail content.
- The detail row's spanning cell remains visible whenever the row is expanded.
- Team changes rows with no event history remain non-expandable.
- Desktop Team changes expansion behavior remains unchanged.
- Focused renderer tests pass, a mobile browser interaction check passes, the checked-in `docs/index.html` matches renderer output, and `ARCHITECTURE.md` describes the delivered behavior.

## Assumptions And Decisions

- This is a focused generated-report bug fix with one renderer and browser-interaction validation path, so Streamlined SDD applies.
- The existing inline `onclick` handler is already activated by a touch-generated click; the issue is presentation, not event delivery.
- The defect occurs because the mobile selector matches the first and only cell in the detail row, even though that cell spans all table columns.
- The corrective selector should scope Manager-cell hiding to headers and non-detail summary rows, rather than changing table markup or interaction logic.
