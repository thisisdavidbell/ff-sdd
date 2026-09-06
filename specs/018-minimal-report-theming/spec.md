---
type: Feature Specification
title: Minimal Report Theming
description: Give the rendered fantasy football report a clearer identity through updated copy, a compact FFD logo, mobile title hierarchy, and a shared accent color.
tags: [report, rendering, theming, logo, mobile, display]
status: draft
feature: 018-minimal-report-theming
sdd_approach: streamlined
input_summary: Rename the report to Fantasy Football Data, add an FFD logo, label its generation time Updated, ensure the mobile title is smaller than section titles, and introduce one basic theme color.
generated: { by: GitHub Copilot, at: 2026-09-06T00:00:00Z }
---

# Requested Behavior

The rendered report's browser title and visible desktop and mobile report titles MUST read `Fantasy Football Data`.

The visible desktop and mobile report headers MUST include a compact, accessible `FFD` lettermark logo. It MUST be implemented in the generated HTML and CSS, without a new image asset or external dependency.

The timestamp label in both report headers MUST read `Updated:`.

On a narrow viewport, the mobile report title MUST render at a larger font size than the report's section headings.

The report MUST use a single shared accent color to provide a recognisable basic visual theme. The accent MUST be applied consistently to the FFD lettermark, primary report identity, and interactive emphasis while preserving the existing light and dark display modes.

The visual treatment MUST remain restrained and information-first, inspired by the clarity and recognisability of BBC and BBC Sport branding rather than a broad redesign.

## Scope Boundaries

- Change only generated report HTML/CSS and focused renderer tests required for this presentation update.
- Preserve report data, tables, chart behavior, navigation, accessibility labels, timestamps, responsive layout, and dark-mode switching.
- Do not change capture, processing, storage formats, external dependencies, or introduce image assets or new configuration.

## Affected Artifacts

- `internal/render/html.go`
- `tests/unit/test_render_output_test.go`
- `docs/index.html`
- `ARCHITECTURE.md`
- `specs/018-minimal-report-theming/spec.md`
- `specs/index.md`
- `specs/log.md`

## Acceptance Checks

- Renderer output includes `Fantasy Football Data` in the document title and both report headers, and excludes the former title.
- Renderer output includes an accessible `FFD` lettermark logo in the desktop and mobile report headers.
- Renderer output labels generated timestamps as `Updated:` in the desktop and mobile headers.
- The narrow-viewport CSS gives `.mobile-title` a font size smaller than `h2`.
- A shared accent-color custom property is defined for both display modes and visibly themes the FFD lettermark and report interaction emphasis.
- The report retains its compact, data-first layout without decorative imagery or a large visual treatment.
- Existing focused renderer tests pass, the checked-in `docs/index.html` matches renderer output, and `ARCHITECTURE.md` describes the delivered presentation behavior.

## Assumptions And Decisions

- This is a focused report display change with one renderer validation path, so Streamlined SDD applies.
- One accent color means a shared semantic color token, with a contrast-appropriate value in each existing display mode, rather than a redesigned palette.
- The FFD logo is a CSS-styled text lettermark so it is sharp, accessible, lightweight, and compatible with static report generation.
- BBC and BBC Sport are visual references for clarity and brand recognition only; no third-party trademarks, typefaces, or visual assets will be reused.
- The title-size requirement applies specifically to the existing narrow-viewport layout and compares the mobile title with `h2` section titles.