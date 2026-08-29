package render

import (
	"bytes"
	"fmt"
	"html"
	"sort"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

// BuildHTML renders a static HTML summary from current ownership and change summaries.
func BuildHTML(current map[string]models.PlayerOwnershipRecord, history map[string][]models.PlayerOwnershipHistoryEntry, summaries []models.ManagerChangeSummary) (string, error) {
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
	b.WriteString("<html><head><meta charset=\"utf-8\"><title>Guy Sports Team Report</title>\n")
	b.WriteString("<style>\n")
	b.WriteString("body{font-family:-apple-system,BlinkMacSystemFont,\"Segoe UI\",Roboto,Helvetica,Arial,sans-serif;max-width:1150px;margin:2rem auto;padding:0 1.5rem;color:#222;line-height:1.5}\n")
	b.WriteString("h1{border-bottom:2px solid #eaeaea;padding-bottom:.5rem;margin-bottom:1.5rem}\n")
	b.WriteString("h2{margin-top:2rem;margin-bottom:.75rem;color:#333}\n")
	b.WriteString("table{border-collapse:collapse;width:100%;margin:1rem 0;box-shadow:0 1px 3px rgba(0,0,0,0.05)}\n")
	b.WriteString("th,td{border:1px solid #ddd;padding:.6rem .8rem;text-align:left}\n")
	b.WriteString("th{background:#f8f9fa;font-weight:600}\n")
	b.WriteString("tr:nth-child(even){background:#fafafa}\n")
	b.WriteString("tr:hover{background:#f1f5f9}\n")
	b.WriteString(".chart-container{width:100%;overflow-x:auto;margin:1.5rem 0;border:1px solid #e2e8f0;border-radius:8px;background:#fff;padding:1rem;box-shadow:0 1px 4px rgba(0,0,0,0.06)}\n")
	b.WriteString(".trend-chart{width:100%;min-height:750px;height:800px;display:block}\n")
	b.WriteString(".player-line{transition:stroke-width 0.2s,opacity 0.2s;opacity:0.85;cursor:pointer}\n")
	b.WriteString(".player-line:hover{stroke-width:4px !important;opacity:1 !important}\n")
	b.WriteString(".player-point{transition:r 0.2s;cursor:pointer}\n")
	b.WriteString(".player-point:hover{r:6}\n")
	b.WriteString(".legend-container{display:flex;flex-wrap:wrap;gap:.5rem;margin-top:1rem;max-height:220px;overflow-y:auto;padding:.5rem;border:1px solid #edf2f7;border-radius:6px;background:#f8fafc}\n")
	b.WriteString(".legend-item{display:inline-flex;align-items:center;font-size:12px;padding:3px 8px;border-radius:4px;background:#fff;border:1px solid #e2e8f0;cursor:pointer;user-select:none;transition:all .15s}\n")
	b.WriteString(".legend-item:hover{border-color:#94a3b8;background:#f1f5f9}\n")
	b.WriteString(".legend-color{width:10px;height:10px;border-radius:50%;margin-right:6px;display:inline-block}\n")
	b.WriteString("</style></head><body>\n")
	b.WriteString("<h1>Guy Sports Team Report</h1>\n")
	b.WriteString("<h2>Player ownership</h2>\n")
	b.WriteString("<table><thead><tr><th>Player</th><th>Managers</th></tr></thead><tbody>\n")
	for _, rec := range playerList {
		name := rec.PlayerName
		if name == "" {
			name = rec.PlayerID
		}
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>\n", html.EscapeString(name), rec.ManagerCount))
	}
	b.WriteString("</tbody></table>\n")

	b.WriteString("<h2>Manager changes</h2>\n")
	b.WriteString("<table><thead><tr><th>Manager</th><th>Team</th><th>Total changes</th><th>Latest change</th></tr></thead><tbody>\n")

	// Sort summaries by ManagerName / TeamName
	sortedSummaries := make([]models.ManagerChangeSummary, len(summaries))
	copy(sortedSummaries, summaries)
	sort.Slice(sortedSummaries, func(i, j int) bool {
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

	for _, summary := range sortedSummaries {
		mgrName := summary.ManagerName
		if mgrName == "" {
			mgrName = summary.ManagerID
		}
		teamName := summary.TeamName
		if teamName == "" {
			teamName = mgrName
		}
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td></tr>\n", html.EscapeString(mgrName), html.EscapeString(teamName), summary.TotalChanges, html.EscapeString(summary.LatestChangeAt)))
	}
	b.WriteString("</tbody></table>\n")

	b.WriteString("<h2>Historical trends</h2>\n")
	renderHistoricalTrendsChart(&b, current, history, playerList)

	b.WriteString("</body></html>\n")
	return b.String(), nil
}

// renderHistoricalTrendsChart builds a large SVG line chart showing player count trends over time with consistent x-axis spacing.
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

	// Compute equidistant X coordinates for timestamps (consistent gaps on x-axis)
	numPoints := len(timestamps)
	xCoords := make([]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		if numPoints == 1 {
			xCoords[i] = leftMargin + plotWidth/2.0
		} else {
			xCoords[i] = leftMargin + float64(i)*(plotWidth/float64(numPoints-1))
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
		label := ts
		if parsed, err := time.Parse(time.RFC3339, ts); err == nil {
			label = parsed.Format("2006-01-02 15:04")
		}
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
			b.WriteString(fmt.Sprintf("<title>%s: %d managers at %s</title>\n", html.EscapeString(player.name), c, html.EscapeString(timestamps[i])))
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
