package render

import (
	"bytes"
	"fmt"
	"html"
	"sort"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

var (
	londonLocation, londonLocationErr = time.LoadLocation("Europe/London")
)

func formatDisplayTimestamp(timestamp string) string {
	parsed, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return timestamp
	}
	return parsed.In(londonLocation).Format("2006-01-02 15:04")
}

// BuildHTML renders a static HTML summary from current ownership and change summaries.
func BuildHTML(current map[string]models.PlayerOwnershipRecord, history map[string][]models.PlayerOwnershipHistoryEntry, summaries []models.ManagerChangeSummary) (string, error) {
	if londonLocationErr != nil {
		return "", fmt.Errorf("load Europe/London location: %w", londonLocationErr)
	}

	now := time.Now().UTC()
	generatedAt := now.In(londonLocation).Format("2006-01-02 15:04")

	// Sort players by manager count descending, then player name ascending
	playerList := make([]models.PlayerOwnershipRecord, 0, len(current))
	for _, rec := range current {
		playerList = append(playerList, rec)
	}
	sort.Slice(playerList, func(i, j int) bool {
		if playerList[i].ManagerCount != playerList[j].ManagerCount {
			return playerList[i].ManagerCount > playerList[j].ManagerCount
		}
		return playerList[i].PlayerName < playerList[j].PlayerName
	})

	var b bytes.Buffer
	b.WriteString("<!doctype html>\n")
	b.WriteString("<html><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>Guy Sports Team Report</title>\n")
	b.WriteString("<script>(function(){try{var theme=localStorage.getItem('report-theme');if(theme==='light'||theme==='dark'){document.documentElement.dataset.theme=theme}else if(window.matchMedia&&window.matchMedia('(prefers-color-scheme: dark)').matches){document.documentElement.dataset.theme='dark'}}catch(e){}}())</script>\n")
	b.WriteString("<style>\n")
	b.WriteString(":root{color-scheme:light;--page-bg:#fff;--text:#222;--heading:#333;--border:#ddd;--subtle-border:#e2e8f0;--surface:#fff;--surface-muted:#f8f9fa;--row:#fafafa;--hover:#f1f5f9;--muted-text:#475569;--chart-bg:#fcfcfd;--grid:#e2e8f0;--grid-light:#f1f5f9;--tick:#94a3b8}\n")
	b.WriteString(":root[data-theme=\"dark\"]{color-scheme:dark;--page-bg:#171a1f;--text:#edf2f7;--heading:#f8fafc;--border:#4a5568;--subtle-border:#3d4757;--surface:#222831;--surface-muted:#2d3642;--row:#202731;--hover:#354151;--muted-text:#cbd5e1;--chart-bg:#1d242d;--grid:#475569;--grid-light:#334155;--tick:#94a3b8}\n")
	b.WriteString("html{scroll-behavior:smooth}body{font-family:-apple-system,BlinkMacSystemFont,\"Segoe UI\",Roboto,Helvetica,Arial,sans-serif;max-width:1150px;margin:2rem auto;padding:0 1.5rem;color:var(--text);background:var(--page-bg);line-height:1.5}\n")
	b.WriteString(".report-header{display:flex;align-items:flex-start;justify-content:space-between;gap:1rem;border-bottom:2px solid var(--subtle-border);padding-bottom:.5rem;margin-bottom:1.5rem}\n")
	b.WriteString("h1{margin:0}.report-generated{margin:.25rem 0 0;color:var(--muted-text);font-size:.875rem}.theme-toggle{flex:none;border:1px solid var(--border);border-radius:4px;background:var(--surface-muted);color:var(--text);padding:.4rem .6rem;cursor:pointer}.theme-toggle:hover{background:var(--hover)}\n")
	b.WriteString("h2{margin-top:2rem;margin-bottom:.75rem;color:var(--heading)}\n")
	b.WriteString("table{border-collapse:collapse;width:100%;margin:1rem 0;box-shadow:0 1px 3px rgba(0,0,0,0.12)}\n")
	b.WriteString("th,td{border:1px solid var(--border);padding:.6rem .8rem;text-align:left}\n")
	b.WriteString("th{background:var(--surface-muted);font-weight:600}\n")
	b.WriteString("tr:nth-child(even){background:var(--row)}\n")
	b.WriteString("tr:hover{background:var(--hover)}\n")
	b.WriteString(".chart-container{width:100%;overflow-x:auto;margin:1.5rem 0;border:1px solid var(--subtle-border);border-radius:8px;background:var(--surface);padding:1rem;box-shadow:0 1px 4px rgba(0,0,0,0.16)}\n")
	b.WriteString(".trend-chart{width:100%;min-height:750px;height:800px;display:block}\n")
	b.WriteString(".player-line{transition:stroke-width 0.2s,opacity 0.2s;opacity:0.85;cursor:pointer}\n")
	b.WriteString(".player-line:hover{stroke-width:4px !important;opacity:1 !important}\n")
	b.WriteString(".player-point{transition:r 0.2s;cursor:pointer}\n")
	b.WriteString(".player-point:hover{r:6}\n")
	b.WriteString(".legend-container{display:flex;flex-wrap:wrap;gap:.5rem;margin-top:1rem;max-height:220px;overflow-y:auto;padding:.5rem;border:1px solid var(--subtle-border);border-radius:6px;background:var(--surface-muted)}\n")
	b.WriteString(".legend-item{display:inline-flex;align-items:center;font-size:12px;padding:3px 8px;border-radius:4px;background:var(--surface);border:1px solid var(--subtle-border);cursor:pointer;user-select:none;transition:all .15s}\n")
	b.WriteString(".legend-item:hover{border-color:var(--tick);background:var(--hover)}\n")
	b.WriteString(".legend-color{width:10px;height:10px;border-radius:50%;margin-right:6px;display:inline-block}\n")
	b.WriteString(".clickable-row{cursor:pointer;user-select:none}\n")
	b.WriteString(".clickable-row:hover{background-color:var(--hover) !important}\n")
	b.WriteString(".manager-detail-row{background-color:var(--surface-muted);display:none}\n")
	b.WriteString(".event-card{margin:.4rem 0;padding:.6rem .8rem;border:1px solid var(--subtle-border);border-radius:6px;background:var(--surface)}\n")
	b.WriteString(".event-header{font-weight:600;font-size:13px;color:var(--muted-text);margin-bottom:.3rem}\n")
	b.WriteString(".added-player{color:#16a34a;font-size:13px;font-weight:500;margin-right:.8rem;display:inline-block}\n")
	b.WriteString(".removed-player{color:#dc2626;font-size:13px;font-weight:500;margin-right:.8rem;display:inline-block}\n")
	b.WriteString(".expand-icon{display:inline-block;font-size:11px;margin-left:6px;color:var(--muted-text)}\n")
	b.WriteString(".date-cell{display:inline-flex;align-items:center;gap:.4rem}.relative-time{font-size:.75rem;color:var(--muted-text)}\n")
	b.WriteString(".pill{display:inline-block;font-size:.7rem;font-weight:700;letter-spacing:.02em;color:#fff;border-radius:999px;padding:.1rem .5rem}.pill-day{background:#2563eb}.pill-week{background:#7c3aed}.pill-month{background:#64748b}\n")
	b.WriteString(".ownership-up{color:#16a34a;font-weight:700}.ownership-down{color:#dc2626;font-weight:700}\n")
	b.WriteString(".table-preview tbody tr.preview-hidden{display:none}.table-preview.is-expanded tbody tr.preview-hidden{display:table-row}.table-expander{margin:-.25rem 0 1rem;border:1px solid var(--border);border-radius:4px;background:var(--surface-muted);color:var(--text);padding:.4rem .6rem;cursor:pointer}.table-expander:hover{background:var(--hover)}\n")
	b.WriteString("@media (max-width:600px){body{margin:1rem auto;padding:0 1rem}.report-header{gap:.75rem}.theme-toggle{padding:.35rem .5rem}}\n")
	b.WriteString(".report-nav{position:sticky;top:0;z-index:10;margin:0 -1.5rem;padding:.5rem 1.5rem;background:var(--surface);border-bottom:1px solid var(--subtle-border)}.report-nav summary{cursor:pointer;font-weight:600;color:var(--heading)}.report-nav summary:focus-visible,.report-nav a:focus-visible{outline:2px solid var(--tick);outline-offset:3px}.report-nav-links{display:grid;gap:.2rem;padding:.5rem 0}.report-nav a{display:block;padding:.45rem;color:var(--text);text-decoration:none;border-left:3px solid transparent}.report-nav a:hover,.report-nav a:focus-visible,.report-nav a[aria-current=\"true\"]{background:var(--hover);border-left-color:var(--tick)}@media(min-width:1000px){body{max-width:1150px;margin:2rem 2rem 2rem 15rem;padding:0 1.5rem}.report-nav{position:fixed;inset:0 auto 0 0;width:12rem;margin:0;padding:2rem 1rem;border-right:1px solid var(--subtle-border);background:var(--surface)}.report-nav summary{display:none}.report-nav:not([open]) .report-nav-links{display:grid}.report-nav-links{padding:0}.report-nav a{border-left:3px solid transparent}}@media(max-width:600px){.report-nav{margin:0 -1rem;padding:.5rem 1rem}}\n")
	b.WriteString("</style></head><body>\n")
	b.WriteString("<header class=\"report-header\"><div><h1>Guy Sports Team Report</h1>\n")
	b.WriteString(fmt.Sprintf("<p class=\"report-generated\">Report generated: %s</p></div><button class=\"theme-toggle\" type=\"button\" aria-label=\"Toggle dark mode\" title=\"Toggle dark mode\">Dark mode</button></header>\n", generatedAt))
	b.WriteString("<nav class=\"report-nav\" aria-label=\"Report sections\"><details open><summary>Report sections</summary><div class=\"report-nav-links\"><a href=\"#top\">Top</a><a href=\"#player-ownership\">Player ownership</a><a href=\"#team-changes\">Team changes</a><a href=\"#historical-trends\">Historical trends</a></div></details></nav>\n")
	b.WriteString("<h2 id=\"player-ownership\" tabindex=\"-1\">Player ownership</h2>\n")
	b.WriteString("<table id=\"player-ownership-table\" class=\"table-preview\"><thead><tr><th>Player</th><th>Managers</th><th>Last change</th></tr></thead><tbody>\n")
	for _, rec := range playerList {
		name := rec.PlayerName
		if name == "" {
			name = rec.PlayerID
		}
		direction, lastChangeAt := latestOwnershipChange(history[rec.PlayerID])
		indicator := ""
		if direction > 0 {
			indicator = " <span class=\"ownership-up\" aria-label=\"Ownership increased\" title=\"Ownership increased\">&#9650;</span>"
		} else if direction < 0 {
			indicator = " <span class=\"ownership-down\" aria-label=\"Ownership decreased\" title=\"Ownership decreased\">&#9660;</span>"
		}
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d%s</td><td>%s</td></tr>\n", html.EscapeString(name), rec.ManagerCount, indicator, html.EscapeString(formatDisplayTimestamp(lastChangeAt))))
	}
	b.WriteString("</tbody></table>\n")
	b.WriteString("<button class=\"table-expander\" type=\"button\" data-table=\"player-ownership-table\">Show all players</button>\n")

	b.WriteString("<h2 id=\"team-changes\" tabindex=\"-1\">Team changes</h2>\n")
	b.WriteString("<table id=\"team-changes-table\" class=\"table-preview\"><thead><tr><th>Manager</th><th>Team</th><th>Total changes</th><th>Latest change</th></tr></thead><tbody>\n")

	// Sort summaries by newest change first, with unchanged managers last.
	sortedSummaries := make([]models.ManagerChangeSummary, len(summaries))
	copy(sortedSummaries, summaries)
	sort.Slice(sortedSummaries, func(i, j int) bool {
		t1, err1 := time.Parse(time.RFC3339, sortedSummaries[i].LatestChangeAt)
		t2, err2 := time.Parse(time.RFC3339, sortedSummaries[j].LatestChangeAt)
		if err1 == nil && err2 == nil && !t1.Equal(t2) {
			return t1.After(t2)
		}
		if err1 == nil && err2 != nil {
			return true
		}
		if err1 != nil && err2 == nil {
			return false
		}
		n1 := sortedSummaries[i].ManagerName
		if n1 == "" {
			n1 = sortedSummaries[i].ManagerID
		}
		n2 := sortedSummaries[j].ManagerName
		if n2 == "" {
			n2 = sortedSummaries[j].ManagerID
		}
		return n1 < n2
	})

	for idx, summary := range sortedSummaries {
		mgrName := summary.ManagerName
		if mgrName == "" {
			mgrName = summary.ManagerID
		}
		teamName := summary.TeamName
		if teamName == "" {
			teamName = mgrName
		}
		detailId := fmt.Sprintf("mgr-detail-%d", idx)

		if len(summary.EventHistory) > 0 {
			b.WriteString(fmt.Sprintf("<tr class=\"clickable-row\" onclick=\"toggleManagerDetail('%s')\"><td>%s</td><td>%s</td><td>%d <span id=\"icon-%s\" class=\"expand-icon\">▶</span></td><td>%s</td></tr>\n",
				detailId,
				html.EscapeString(mgrName),
				html.EscapeString(teamName),
				summary.TotalChanges,
				detailId,
				formatLatestChangeCell(summary.LatestChangeAt, now),
			))

			// Render detail row
			b.WriteString(fmt.Sprintf("<tr id=\"%s\" class=\"manager-detail-row\"><td colspan=\"4\">\n", detailId))
			for _, event := range summary.EventHistory {
				formattedTime := formatDisplayTimestamp(event.ToCapturedAt)
				b.WriteString("<div class=\"event-card\">\n")
				b.WriteString(fmt.Sprintf("<div class=\"event-header\">Change at %s (%d change%s)</div>\n",
					html.EscapeString(formattedTime),
					event.ChangeCount,
					pluralSuffix(event.ChangeCount),
				))
				b.WriteString("<div>\n")
				for _, p := range event.AddedPlayers {
					pName := p.Name
					if pName == "" {
						pName = p.PlayerID
					}
					b.WriteString(fmt.Sprintf("<span class=\"added-player\">+ %s</span>\n", html.EscapeString(pName)))
				}
				for _, p := range event.RemovedPlayers {
					pName := p.Name
					if pName == "" {
						pName = p.PlayerID
					}
					b.WriteString(fmt.Sprintf("<span class=\"removed-player\">- %s</span>\n", html.EscapeString(pName)))
				}
				b.WriteString("</div>\n")
				b.WriteString("</div>\n")
			}
			b.WriteString("</td></tr>\n")
		} else {
			b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td></tr>\n", html.EscapeString(mgrName), html.EscapeString(teamName), summary.TotalChanges, formatLatestChangeCell(summary.LatestChangeAt, now)))
		}
	}
	b.WriteString("</tbody></table>\n")
	b.WriteString("<button class=\"table-expander\" type=\"button\" data-table=\"team-changes-table\">Show all teams</button>\n")

	b.WriteString("<h2 id=\"historical-trends\" tabindex=\"-1\">Historical trends</h2>\n")
	renderHistoricalTrendsChart(&b, current, history, playerList)

	b.WriteString("<script>\n")
	b.WriteString("function toggleManagerDetail(id){\n")
	b.WriteString("  var row = document.getElementById(id);\n")
	b.WriteString("  var icon = document.getElementById('icon-' + id);\n")
	b.WriteString("  if(!row) return;\n")
	b.WriteString("  if(row.style.display === 'none' || row.style.display === ''){\n")
	b.WriteString("    row.style.display = 'table-row';\n")
	b.WriteString("    if(icon) icon.textContent = '▼';\n")
	b.WriteString("  }else{\n")
	b.WriteString("    row.style.display = 'none';\n")
	b.WriteString("    if(icon) icon.textContent = '▶';\n")
	b.WriteString("  }\n")
	b.WriteString("}\n")
	b.WriteString("(function(){var button=document.querySelector('.theme-toggle');if(!button)return;function update(theme){document.documentElement.dataset.theme=theme;button.textContent=theme==='dark'?'Light mode':'Dark mode';button.setAttribute('aria-label',theme==='dark'?'Switch to light mode':'Switch to dark mode');button.title=button.getAttribute('aria-label')}var current=document.documentElement.dataset.theme==='dark'?'dark':'light';update(current);button.addEventListener('click',function(){var next=document.documentElement.dataset.theme==='dark'?'light':'dark';update(next);try{localStorage.setItem('report-theme',next)}catch(e){}})}())\n")
	b.WriteString("(function(){Array.prototype.forEach.call(document.querySelectorAll('.table-expander'),function(button){var table=document.getElementById(button.dataset.table);if(!table)return;var rows=Array.prototype.filter.call(table.tBodies[0].rows,function(row){return !row.classList.contains('manager-detail-row')});var hiddenRows=rows.slice(10);var itemLabel=button.dataset.table==='player-ownership-table'?'players':'teams';if(!hiddenRows.length){button.hidden=true;return}hiddenRows.forEach(function(row){row.classList.add('preview-hidden')});button.textContent='Show all '+rows.length+' '+itemLabel;button.setAttribute('aria-expanded','false');button.setAttribute('aria-controls',table.id);button.addEventListener('click',function(){var expanded=table.classList.toggle('is-expanded');button.textContent=expanded?'Show fewer '+itemLabel:'Show all '+rows.length+' '+itemLabel;button.setAttribute('aria-expanded',String(expanded))})})}())\n")
	b.WriteString("</script>\n")
	b.WriteString("<script>(function(){var links=document.querySelectorAll('.report-nav-links a[href^=\"#\"]');var sections=Array.prototype.map.call(links,function(link){var target=link.getAttribute('href');return target==='#top'?null:document.querySelector(target)}).filter(Boolean);function update(){var current=sections[0];sections.forEach(function(section){if(section.getBoundingClientRect().top<=120){current=section}});links.forEach(function(link){var target=link.getAttribute('href');link.setAttribute('aria-current',target!=='#top'&&document.querySelector(target)===current?'true':'false')})}window.addEventListener('scroll',update,{passive:true});window.addEventListener('hashchange',update);update()}())</script>\n")

	b.WriteString("</body></html>\n")
	return b.String(), nil
}

func pluralSuffix(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// latestOwnershipChange returns the direction and timestamp of the last count change.
func latestOwnershipChange(entries []models.PlayerOwnershipHistoryEntry) (int, string) {
	if len(entries) < 2 {
		return 0, ""
	}

	direction := 0
	lastChangeAt := ""
	previousCount := entries[0].ManagerCount
	for _, entry := range entries[1:] {
		if entry.ManagerCount == previousCount {
			continue
		}
		if entry.ManagerCount > previousCount {
			direction = 1
		} else {
			direction = -1
		}
		lastChangeAt = entry.CapturedAt
		previousCount = entry.ManagerCount
	}
	return direction, lastChangeAt
}

// formatLatestChangeCell renders the "Latest change" cell: the raw value plus a
// relative-time phrase and, for changes within the last month, a single
// Day/Week/Month tiered pill (see specs/005-recent-manager-changes-highlight).
func formatLatestChangeCell(latestChangeAt string, now time.Time) string {
	escaped := html.EscapeString(formatDisplayTimestamp(latestChangeAt))
	if latestChangeAt == "" {
		return escaped
	}
	t, err := time.Parse(time.RFC3339, latestChangeAt)
	if err != nil {
		return escaped
	}
	age := now.Sub(t)
	cell := fmt.Sprintf("<span class=\"date-cell\">%s <span class=\"relative-time\">(%s)</span>", escaped, html.EscapeString(relativeTime(age)))
	if class, label := recencyTier(age); class != "" {
		cell += fmt.Sprintf(" <span class=\"pill %s\">%s</span>", class, label)
	}
	cell += "</span>"
	return cell
}

// relativeTime formats a duration since report generation as a human-readable phrase.
func relativeTime(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d / time.Minute)
		return fmt.Sprintf("%d minute%s", mins, pluralSuffix(mins))
	case d < 24*time.Hour:
		hours := int(d / time.Hour)
		return fmt.Sprintf("%d hour%s", hours, pluralSuffix(hours))
	case d < 30*24*time.Hour:
		days := int(d / (24 * time.Hour))
		return fmt.Sprintf("%d day%s", days, pluralSuffix(days))
	case d < 365*24*time.Hour:
		months := int(d / (30 * 24 * time.Hour))
		return fmt.Sprintf("%d month%s", months, pluralSuffix(months))
	default:
		years := int(d / (365 * 24 * time.Hour))
		return fmt.Sprintf("%d year%s", years, pluralSuffix(years))
	}
}

// recencyTier returns the CSS class and label for the most specific Day/Week/Month
// tier the duration falls into, or empty strings if older than a month.
func recencyTier(d time.Duration) (class string, label string) {
	switch {
	case d < 24*time.Hour:
		return "pill-day", "Day"
	case d < 7*24*time.Hour:
		return "pill-week", "Week"
	case d < 30*24*time.Hour:
		return "pill-month", "Month"
	default:
		return "", ""
	}
}

// renderHistoricalTrendsChart builds a large SVG line chart showing player count trends over elapsed time.
func renderHistoricalTrendsChart(b *bytes.Buffer, current map[string]models.PlayerOwnershipRecord, history map[string][]models.PlayerOwnershipHistoryEntry, playerList []models.PlayerOwnershipRecord) {
	// Collect all unique capture timestamps
	timestampSet := make(map[string]bool)
	maxCount := 0
	for _, entries := range history {
		for _, e := range entries {
			if e.CapturedAt != "" {
				timestampSet[e.CapturedAt] = true
				if e.ManagerCount > maxCount {
					maxCount = e.ManagerCount
				}
			}
		}
	}
	for _, p := range current {
		if p.ManagerCount > maxCount {
			maxCount = p.ManagerCount
		}
	}

	timestamps := make([]string, 0, len(timestampSet))
	for ts := range timestampSet {
		timestamps = append(timestamps, ts)
	}

	// Sort timestamps chronologically
	sort.Slice(timestamps, func(i, j int) bool {
		t1, err1 := time.Parse(time.RFC3339, timestamps[i])
		t2, err2 := time.Parse(time.RFC3339, timestamps[j])
		if err1 != nil || err2 != nil {
			return timestamps[i] < timestamps[j]
		}
		return t1.Before(t2)
	})

	if len(timestamps) == 0 {
		b.WriteString("<p>No historical trend data available.</p>\n")
		return
	}

	// Chart dimensions
	const svgWidth = 1100.0
	const svgHeight = 800.0
	const leftMargin = 60.0
	const rightMargin = 200.0
	const topMargin = 50.0
	const bottomMargin = 100.0
	const plotWidth = svgWidth - leftMargin - rightMargin
	const plotHeight = svgHeight - topMargin - bottomMargin

	// Determine Y-axis max and grid steps
	if maxCount < 10 {
		maxCount = 10
	}
	yMax := ((maxCount + 4) / 5) * 5
	yStep := 5
	if yMax <= 10 {
		yStep = 2
	} else if yMax > 40 {
		yStep = 10
	}

	// Position each timestamp proportionally within the exact capture-time range.
	numPoints := len(timestamps)
	xCoords := make([]float64, numPoints)
	if numPoints == 1 {
		xCoords[0] = leftMargin + plotWidth/2.0
	} else {
		firstTime, _ := time.Parse(time.RFC3339, timestamps[0])
		lastTime, _ := time.Parse(time.RFC3339, timestamps[numPoints-1])
		duration := lastTime.Sub(firstTime)
		for i, timestamp := range timestamps {
			capturedAt, _ := time.Parse(time.RFC3339, timestamp)
			if duration == 0 {
				xCoords[i] = leftMargin + plotWidth/2.0
				continue
			}
			xCoords[i] = leftMargin + (float64(capturedAt.Sub(firstTime)) / float64(duration)) * plotWidth
		}
	}

	b.WriteString("<div class=\"chart-container\">\n")
	b.WriteString(fmt.Sprintf("<svg class=\"trend-chart\" viewBox=\"0 0 %d %d\" preserveAspectRatio=\"xMidYMid meet\">\n", int(svgWidth), int(svgHeight)))

	// Chart background
	b.WriteString(fmt.Sprintf("<rect x=\"%.1f\" y=\"%.1f\" width=\"%.1f\" height=\"%.1f\" fill=\"#fcfcfd\" stroke=\"#e2e8f0\" stroke-width=\"1\" />\n", leftMargin, topMargin, plotWidth, plotHeight))

	// Horizontal grid lines and Y-axis labels
	for v := 0; v <= yMax; v += yStep {
		y := topMargin + plotHeight - (float64(v)/float64(yMax))*plotHeight
		b.WriteString(fmt.Sprintf("<line x1=\"%.1f\" y1=\"%.1f\" x2=\"%.1f\" y2=\"%.1f\" stroke=\"#e2e8f0\" stroke-dasharray=\"4,4\" stroke-width=\"1\" />\n", leftMargin, y, leftMargin+plotWidth, y))
		b.WriteString(fmt.Sprintf("<text x=\"%.1f\" y=\"%.1f\" text-anchor=\"end\" font-size=\"12\" fill=\"#64748b\" font-family=\"sans-serif\">%d</text>\n", leftMargin-8, y+4, v))
	}

	// Y-axis label title
	b.WriteString(fmt.Sprintf("<text x=\"%.1f\" y=\"%.1f\" text-anchor=\"middle\" transform=\"rotate(-90 %.1f, %.1f)\" font-size=\"13\" font-weight=\"600\" fill=\"#475569\" font-family=\"sans-serif\">Manager Count</text>\n", 18.0, topMargin+plotHeight/2.0, 18.0, topMargin+plotHeight/2.0))

	// Vertical grid lines, ticks, and X-axis date labels
	for i, ts := range timestamps {
		x := xCoords[i]
		// Vertical grid line
		b.WriteString(fmt.Sprintf("<line x1=\"%.1f\" y1=\"%.1f\" x2=\"%.1f\" y2=\"%.1f\" stroke=\"#f1f5f9\" stroke-width=\"1\" />\n", x, topMargin, x, topMargin+plotHeight))
		// Tick mark
		b.WriteString(fmt.Sprintf("<line x1=\"%.1f\" y1=\"%.1f\" x2=\"%.1f\" y2=\"%.1f\" stroke=\"#94a3b8\" stroke-width=\"1.5\" />\n", x, topMargin+plotHeight, x, topMargin+plotHeight+6))

		// Formatted date label
		label := formatDisplayTimestamp(ts)
		b.WriteString(fmt.Sprintf("<text x=\"%.1f\" y=\"%.1f\" transform=\"rotate(-35 %.1f, %.1f)\" text-anchor=\"end\" font-size=\"11\" fill=\"#475569\" font-family=\"sans-serif\">%s</text>\n", x, topMargin+plotHeight+20, x, topMargin+plotHeight+20, html.EscapeString(label)))
	}

	// Determine active players (from playerList or history)
	type playerChartInfo struct {
		id          string
		name        string
		color       string
		counts      []int
		lastCount   int
		playerIndex int
	}

	// Track all player keys
	playerKeysSet := make(map[string]bool)
	for _, p := range playerList {
		playerKeysSet[p.PlayerID] = true
	}
	for id := range history {
		playerKeysSet[id] = true
	}

	orderedPlayers := make([]playerChartInfo, 0, len(playerKeysSet))
	idx := 0
	for _, p := range playerList {
		pID := p.PlayerID
		pName := p.PlayerName
		if pName == "" {
			pName = pID
		}

		countMap := make(map[string]int)
		for _, e := range history[pID] {
			countMap[e.CapturedAt] = e.ManagerCount
		}

		counts := make([]int, numPoints)
		lastC := 0
		for tIdx, ts := range timestamps {
			c := countMap[ts]
			counts[tIdx] = c
			lastC = c
		}
		if lastC == 0 && p.ManagerCount > 0 {
			lastC = p.ManagerCount
		}

		hue := int(float64(idx)*137.508) % 360
		color := fmt.Sprintf("hsl(%d, 70%%, 42%%)", hue)

		orderedPlayers = append(orderedPlayers, playerChartInfo{
			id:          pID,
			name:        pName,
			color:       color,
			counts:      counts,
			lastCount:   lastC,
			playerIndex: idx,
		})
		delete(playerKeysSet, pID)
		idx++
	}

	// Any remaining historical players not in current playerList
	for pID := range playerKeysSet {
		pName := pID
		countMap := make(map[string]int)
		for _, e := range history[pID] {
			countMap[e.CapturedAt] = e.ManagerCount
		}
		counts := make([]int, numPoints)
		lastC := 0
		for tIdx, ts := range timestamps {
			c := countMap[ts]
			counts[tIdx] = c
			lastC = c
		}
		hue := int(float64(idx)*137.508) % 360
		color := fmt.Sprintf("hsl(%d, 70%%, 42%%)", hue)
		orderedPlayers = append(orderedPlayers, playerChartInfo{
			id:          pID,
			name:        pName,
			color:       color,
			counts:      counts,
			lastCount:   lastC,
			playerIndex: idx,
		})
		idx++
	}

	// Render player lines and data points
	for _, player := range orderedPlayers {
		var pathBuf bytes.Buffer
		hasPositiveCount := false
		for i, c := range player.counts {
			if c > 0 {
				hasPositiveCount = true
			}
			y := topMargin + plotHeight - (float64(c)/float64(yMax))*plotHeight
			if i == 0 {
				pathBuf.WriteString(fmt.Sprintf("M %.1f %.1f", xCoords[i], y))
			} else {
				pathBuf.WriteString(fmt.Sprintf(" L %.1f %.1f", xCoords[i], y))
			}
		}

		if numPoints > 1 {
			b.WriteString(fmt.Sprintf("<path id=\"line-%s\" class=\"player-line\" data-player-id=\"%s\" d=\"%s\" stroke=\"%s\" stroke-width=\"2\" fill=\"none\">\n", html.EscapeString(player.id), html.EscapeString(player.id), pathBuf.String(), player.color))
			b.WriteString(fmt.Sprintf("<title>%s (Current: %d)</title>\n", html.EscapeString(player.name), player.lastCount))
			b.WriteString("</path>\n")
		}

		// Render points on the line
		for i, c := range player.counts {
			y := topMargin + plotHeight - (float64(c)/float64(yMax))*plotHeight
			b.WriteString(fmt.Sprintf("<circle cx=\"%.1f\" cy=\"%.1f\" r=\"3.5\" fill=\"%s\" class=\"player-point\" data-player-id=\"%s\">\n", xCoords[i], y, player.color, html.EscapeString(player.id)))
			b.WriteString(fmt.Sprintf("<title>%s: %d managers at %s</title>\n", html.EscapeString(player.name), c, html.EscapeString(formatDisplayTimestamp(timestamps[i]))))
			b.WriteString("</circle>\n")
		}

		// End-of-line label for visible players
		if hasPositiveCount && player.lastCount > 0 {
			lastX := xCoords[numPoints-1]
			lastY := topMargin + plotHeight - (float64(player.lastCount)/float64(yMax))*plotHeight
			b.WriteString(fmt.Sprintf("<text x=\"%.1f\" y=\"%.1f\" font-size=\"11\" font-weight=\"500\" fill=\"%s\" font-family=\"sans-serif\" data-player-id=\"%s\">%s (%d)</text>\n", lastX+8, lastY+3.5, player.color, html.EscapeString(player.id), html.EscapeString(player.name), player.lastCount))
		}
	}

	b.WriteString("</svg>\n")
	b.WriteString("</div>\n")

	// Legend with player badges and interactive highlighting
	b.WriteString("<div class=\"legend-container\" id=\"trend-legend\">\n")
	for _, player := range orderedPlayers {
		b.WriteString(fmt.Sprintf("<span class=\"legend-item\" data-player-id=\"%s\" onmouseover=\"highlightPlayer('%s')\" onmouseout=\"resetHighlight()\">\n", html.EscapeString(player.id), html.EscapeString(player.id)))
		b.WriteString(fmt.Sprintf("<span class=\"legend-color\" style=\"background-color:%s;\"></span>", player.color))
		b.WriteString(fmt.Sprintf("<strong>%s</strong>&nbsp;(%d)", html.EscapeString(player.name), player.lastCount))
		b.WriteString("</span>\n")
	}
	b.WriteString("</div>\n")

	// Interactive JavaScript for smooth player inspection
	b.WriteString("<script>\n")
	b.WriteString("function highlightPlayer(id) {\n")
	b.WriteString("  document.querySelectorAll('.player-line').forEach(function(el) {\n")
	b.WriteString("    if (el.getAttribute('data-player-id') === id) {\n")
	b.WriteString("      el.style.strokeWidth = '4.5px';\n")
	b.WriteString("      el.style.opacity = '1';\n")
	b.WriteString("    } else {\n")
	b.WriteString("      el.style.opacity = '0.15';\n")
	b.WriteString("    }\n")
	b.WriteString("  });\n")
	b.WriteString("}\n")
	b.WriteString("function resetHighlight() {\n")
	b.WriteString("  document.querySelectorAll('.player-line').forEach(function(el) {\n")
	b.WriteString("    el.style.strokeWidth = '2px';\n")
	b.WriteString("    el.style.opacity = '0.85';\n")
	b.WriteString("  });\n")
	b.WriteString("}\n")
	b.WriteString("</script>\n")
}

// BuildHTMLFromDoc renders static HTML directly from a PlayerOwnershipDoc and change summaries.
func BuildHTMLFromDoc(doc models.PlayerOwnershipDoc, summaries []models.ManagerChangeSummary) (string, error) {
	current := make(map[string]models.PlayerOwnershipRecord, len(doc.Players))
	history := make(map[string][]models.PlayerOwnershipHistoryEntry, len(doc.Players))

	for _, p := range doc.Players {
		current[p.PlayerID] = models.PlayerOwnershipRecord{
			PlayerID:     p.PlayerID,
			PlayerName:   p.PlayerName,
			ManagerCount: p.CurrentCount,
		}
		history[p.PlayerID] = p.History
	}

	return BuildHTML(current, history, summaries)
}
