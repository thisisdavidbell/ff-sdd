---
type: Feature Specification
title: Minimal Report Theming
description: Give the rendered fantasy football report a clearer identity through updated copy, a compact FFD logo, mobile title hierarchy, and a shared accent color.
tags: [report, rendering, theming, logo, mobile, display]
status: stable
feature: 018-minimal-report-theming
sdd_approach: streamlined
input_summary: Rename the report to Fantasy Football Data, add an FFD logo, label its generation time Updated, ensure the mobile title is larger than section titles, and introduce one basic theme color.
generated: { by: GitHub Copilot, at: 2026-09-06T00:00:00Z }
---

# Requested Behavior

The rendered report's browser title and visible desktop and mobile report titles MUST read `Fantasy Football Data`.

The mobile report header MUST include a compact, accessible `FFD` lettermark logo. On desktop viewports, the selected FFD logo MUST appear above the sidebar navigation menu rather than in the main page header.

The desktop sidebar logo MUST render at `9rem` wide by `5.4rem` high, and the navigation content MUST begin below it. The main desktop page header MUST retain the report title and timestamp without a logo.

The timestamp label in both report headers MUST read `Updated:`.

On a narrow viewport, the mobile report title MUST render at a larger font size than the report's section headings. The mobile FFD logo MUST render at `5rem` wide by `3rem` high.

The report MUST use a single shared accent color to provide a recognisable basic visual theme. The accent MUST be applied consistently to the FFD lettermark, primary report identity, and interactive emphasis while preserving the existing light and dark display modes.

The desktop and mobile long-form report titles MUST use the selected option B treatment: all of `Fantasy Football Data` in the shared blue accent color.

The visual treatment MUST remain restrained and information-first, inspired by the clarity and recognisability of BBC and BBC Sport branding rather than a broad redesign.

## Scope Boundaries

- Create a static comparison prototype under `examples/` before selecting the generated report HTML/CSS treatment.
- Change only generated report HTML/CSS and focused renderer tests required for the selected presentation update.
- Preserve report data, tables, chart behavior, navigation, accessibility labels, timestamps, responsive layout, and dark-mode switching.
- Do not change capture, processing, storage formats, external dependencies, or introduce new configuration.

## Affected Artifacts

- `internal/render/html.go`
- `tests/unit/test_render_output_test.go`
- `docs/index.html`
- `docs/assets/ffd-logo.svg`
- `docs/assets/ffd-favicon.svg`
- `ARCHITECTURE.md`
- `specs/018-minimal-report-theming/spec.md`
- `examples/018-minimal-report-theming/index.html`
- `specs/index.md`
- `specs/log.md`

## Acceptance Checks

- Renderer output includes `Fantasy Football Data` in the document title and both report headers, and excludes the former title.
- Renderer output includes an accessible `FFD` lettermark logo above desktop sidebar navigation and in the mobile report header.
- Renderer output labels generated timestamps as `Updated:` in the desktop and mobile headers.
- The narrow-viewport CSS gives `.mobile-title` a font size larger than `h2`.
- Narrow-viewport CSS renders the mobile FFD logo at `5rem` wide and `3rem` high.
- Desktop CSS renders the sidebar logo at `9rem` wide and `5.4rem` high, with navigation content below it.
- A shared accent-color custom property is defined for both display modes and visibly themes the FFD lettermark and report interaction emphasis.
- The desktop and mobile long-form report titles use the shared accent color.
- The report retains its compact, data-first layout without decorative imagery or a large visual treatment.
- A local comparison prototype presents multiple restrained accent and FFD lettermark directions for review before one is applied to the report.
- The local comparison prototype offers a toggle to review every direction in light and dark modes.
- The local comparison prototype presents blue long-form title treatments, including initial letters, all words, Fantasy Football, Data, a rule, and a bar, for user selection.
- The local comparison prototype includes separate full-report review pages for title option A (blue initial letters) and title option B (all-blue title).
- Existing focused renderer tests pass, the checked-in `docs/index.html` matches renderer output, and `ARCHITECTURE.md` describes the delivered presentation behavior.

## Assumptions And Decisions

- This is a focused report display change with one renderer validation path, so Streamlined SDD applies.
- One accent color means a shared semantic color token, with a contrast-appropriate value in each existing display mode, rather than a redesigned palette.
- The FFD logo is a CSS-styled text lettermark so it is sharp, accessible, lightweight, and compatible with static report generation.
- The selected direction is 3a: a blue framed FFD lettermark with a solid blue circle at the lower right. It is stored as a standalone SVG and used by the report headers and favicon.
- The favicon uses a square, small-size-optimized rendering of the selected FFD identity so the framed letters and solid-circle detail remain legible in a browser tab.
- BBC and BBC Sport are visual references for clarity and brand recognition only; no third-party trademarks, typefaces, or visual assets will be reused.
- The title-size requirement applies specifically to the existing narrow-viewport layout and compares the mobile title with `h2` section titles.
- The enlarged logo and moved navigation apply only to the desktop sidebar; the compact logo remains in the mobile header. The `9rem` by `5.4rem` desktop dimensions are a fifty-percent increase from the prior sidebar size.
- The long-form title-color treatments are exploratory prototype work; no title treatment changes the generated report until the user selects one.
- The user selected option B, which colors the entire long-form title blue.
- The full-report title option pages load the generated report and apply only their selected title treatment, so report content and interaction remain representative.