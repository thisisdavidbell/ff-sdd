package unit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/render"
)

func TestRenderHTMLIncludesOwnershipAndChanges(t *testing.T) {
	current := map[string]models.PlayerOwnershipRecord{
		"p1": {PlayerID: "p1", PlayerName: "Haaland", ManagerCount: 5},
		"p2": {PlayerID: "p2", PlayerName: "Salah", ManagerCount: 10},
		"p3": {PlayerID: "p3", PlayerName: "Saka", ManagerCount: 2},
	}
	history := map[string][]models.PlayerOwnershipHistoryEntry{
		"p1": {{CapturedAt: "2026-01-10T12:00:00Z", ManagerCount: 4}, {CapturedAt: "2026-01-10T14:00:00Z", ManagerCount: 5}},
		"p2": {{CapturedAt: "2026-01-10T12:00:00Z", ManagerCount: 8}, {CapturedAt: "2026-01-10T14:00:00Z", ManagerCount: 10}},
	}
	summaries := []models.ManagerChangeSummary{
		{
			ManagerID:                "8",
			ManagerName:              "John_Doe",
			TeamName:                 "Finsbury_Bridge",
			TotalChanges:             2,
			LatestChangeAt:           "2026-01-10T15:00:00Z",
			ChangedSinceLastSnapshot: true,
		},
	}

	html, err := render.BuildHTML(current, history, summaries)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	// Verify player table exists and is sorted descending (Salah with 10 before Haaland with 5 before Saka with 2)
	salahPos := strings.Index(html, "Salah")
	haalandPos := strings.Index(html, "Haaland")
	sakaPos := strings.Index(html, "Saka")

	if salahPos == -1 || haalandPos == -1 || sakaPos == -1 {
		t.Fatalf("expected all players in HTML, got: %s", html)
	}
	if !(salahPos < haalandPos && haalandPos < sakaPos) {
		t.Fatalf("expected descending order (Salah > Haaland > Saka), got positions: Salah=%d, Haaland=%d, Saka=%d", salahPos, haalandPos, sakaPos)
	}

	// Verify manager changes table shows Manager Name and Team Name
	if !strings.Contains(html, "John_Doe") || !strings.Contains(html, "Finsbury_Bridge") {
		t.Fatalf("manager changes missing manager/team name in HTML: %s", html)
	}
	if !strings.Contains(html, "<th>Manager</th><th>Team</th>") {
		t.Fatalf("manager changes table headers missing Manager/Team columns: %s", html)
	}

	// Verify large historical trends line chart is rendered
	if !strings.Contains(html, "<svg class=\"trend-chart\"") {
		t.Fatalf("historical trends line chart SVG missing from HTML: %s", html)
	}
	if !strings.Contains(html, `<meta name="viewport" content="width=device-width, initial-scale=1">`) {
		t.Fatalf("expected standard responsive viewport metadata: %s", html)
	}
	if !strings.Contains(html, "viewBox=\"0 0 1100 800\"") {
		t.Fatalf("expected chart viewBox in HTML: %s", html)
	}
	if !strings.Contains(html, ".trend-chart{min-height:0;height:auto;aspect-ratio:11/8}") && !strings.Contains(html, ".trend-chart{min-height:0;height:auto;aspect-ratio:11 / 8}") {
		t.Fatalf("expected responsive mobile chart height: %s", html)
	}
	if !strings.Contains(html, "class=\"player-line\"") {
		t.Fatalf("expected player line paths in HTML: %s", html)
	}
}

func TestRenderHTMLIncludesReportNavigation(t *testing.T) {
	html, err := render.BuildHTML(nil, nil, nil)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	expectedNavigation := `<span class="nav-label">Statistics</span><a href="#player-ownership">Player ownership</a><a href="#team-changes">Team changes</a></div><div class="nav-group"><span class="nav-label">Trends</span><a href="#historical-trends">Historical trends</a>`
	if !strings.Contains(html, expectedNavigation) {
		t.Fatalf("expected grouped report navigation, got: %s", html)
	}
	if strings.Contains(html, `href="#top"`) {
		t.Fatalf("expected Top navigation link to be removed: %s", html)
	}
	for _, section := range []string{
		`<h2 id="player-ownership" tabindex="-1">Player ownership</h2>`,
		`<h2 id="team-changes" tabindex="-1">Team changes</h2>`,
		`<h2 id="historical-trends" tabindex="-1">Historical trends</h2>`,
	} {
		if !strings.Contains(html, section) {
			t.Fatalf("expected anchored report section %q", section)
		}
	}
	if !strings.Contains(html, `class="report-nav"`) || !strings.Contains(html, `class="mobile-header"`) || !strings.Contains(html, `class="menu-button"`) || !strings.Contains(html, `aria-expanded="false"`) {
		t.Fatalf("expected responsive report sidebar navigation: %s", html)
	}
}

func TestRenderHTMLShowsLatestOwnershipDirectionAndChangeTime(t *testing.T) {

	current := map[string]models.PlayerOwnershipRecord{
		"up":      {PlayerID: "up", PlayerName: "Rising", ManagerCount: 4},
		"down":    {PlayerID: "down", PlayerName: "Falling", ManagerCount: 2},
		"steady":  {PlayerID: "steady", PlayerName: "Steady", ManagerCount: 3},
		"single":  {PlayerID: "single", PlayerName: "New", ManagerCount: 1},
	}
	history := map[string][]models.PlayerOwnershipHistoryEntry{
		"up":     {{CapturedAt: "2026-01-01T12:00:00Z", ManagerCount: 2}, {CapturedAt: "2026-01-02T12:00:00Z", ManagerCount: 2}, {CapturedAt: "2026-01-03T12:00:00Z", ManagerCount: 4}},
		"down":   {{CapturedAt: "2026-01-01T12:00:00Z", ManagerCount: 5}, {CapturedAt: "2026-01-03T12:00:00Z", ManagerCount: 2}},
		"steady": {{CapturedAt: "2026-01-01T12:00:00Z", ManagerCount: 3}, {CapturedAt: "2026-01-03T12:00:00Z", ManagerCount: 3}},
		"single": {{CapturedAt: "2026-01-03T12:00:00Z", ManagerCount: 1}},
	}

	html, err := render.BuildHTML(current, history, nil)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	if !strings.Contains(html, "<th>Last change</th>") {
		t.Fatalf("expected a Last change ownership column: %s", html)
	}
	if !strings.Contains(html, "Rising</td><td>4 <span class=\"ownership-up\"") || !strings.Contains(html, "2026-01-03 12:00") {
		t.Fatalf("expected latest increase indicator and timestamp: %s", html)
	}
	if !strings.Contains(html, "Falling</td><td>2 <span class=\"ownership-down\"") {
		t.Fatalf("expected latest decrease indicator: %s", html)
	}
	if strings.Contains(html, "Steady</td><td>3 <span class=\"ownership-") || strings.Contains(html, "New</td><td>1 <span class=\"ownership-") {
		t.Fatalf("expected no direction indicator for unchanged ownership: %s", html)
	}
}

func TestRenderHTMLOrdersManagerChangesByLatestUpdate(t *testing.T) {
	summaries := []models.ManagerChangeSummary{
		{ManagerID: "none", ManagerName: "No changes", TeamName: "No changes", TotalChanges: 0},
		{ManagerID: "older", ManagerName: "Older", TeamName: "Older", TotalChanges: 1, LatestChangeAt: "2026-01-01T12:00:00Z"},
		{ManagerID: "newer", ManagerName: "Newer", TeamName: "Newer", TotalChanges: 1, LatestChangeAt: "2026-01-03T12:00:00Z"},
	}

	html, err := render.BuildHTML(map[string]models.PlayerOwnershipRecord{}, nil, summaries)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	newerPos := strings.Index(html, "Newer</td>")
	olderPos := strings.Index(html, "Older</td>")
	nonePos := strings.Index(html, "No changes</td>")
	if !(newerPos < olderPos && olderPos < nonePos) {
		t.Fatalf("expected manager changes ordered newest to oldest with no changes last, got Newer=%d Older=%d No changes=%d", newerPos, olderPos, nonePos)
	}
}

func TestRenderHTMLClassifiesLatestChangeRelativeTimePills(t *testing.T) {
	now := time.Now().UTC()
	summaries := []models.ManagerChangeSummary{
		{ManagerID: "hours", ManagerName: "Hours", TeamName: "Hours", TotalChanges: 1, LatestChangeAt: now.Add(-2 * time.Hour).Format(time.RFC3339)},
		{ManagerID: "days", ManagerName: "Days", TeamName: "Days", TotalChanges: 1, LatestChangeAt: now.Add(-2 * 24 * time.Hour).Format(time.RFC3339)},
		{ManagerID: "weeks", ManagerName: "Weeks", TeamName: "Weeks", TotalChanges: 1, LatestChangeAt: now.Add(-14 * 24 * time.Hour).Format(time.RFC3339)},
		{ManagerID: "months", ManagerName: "Months", TeamName: "Months", TotalChanges: 1, LatestChangeAt: now.Add(-60 * 24 * time.Hour).Format(time.RFC3339)},
	}

	html, err := render.BuildHTML(nil, nil, summaries)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	for _, expected := range []string{
		`<span class="relative-time">2 hours ago</span> <span class="pill pill-hours">Hours</span>`,
		`<span class="relative-time">2 days ago</span> <span class="pill pill-days">Days</span>`,
		`<span class="relative-time">2 weeks ago</span> <span class="pill pill-weeks">Weeks</span>`,
		`<span class="relative-time">2 months ago</span> <span class="pill pill-months">Months</span>`,
	} {
		if !strings.Contains(html, expected) {
			t.Fatalf("expected latest-change presentation %q: %s", expected, html)
		}
	}
	if !strings.Contains(html, `.pill-hours{background:#c026d3}`) {
		t.Fatalf("expected permanent fuchsia Hours pill style: %s", html)
	}
	if strings.Contains(html, `hours-pill-prototype`) || strings.Contains(html, `data-hours-pill-color`) {
		t.Fatalf("expected temporary Hours-pill prototype control to be removed: %s", html)
	}

	tableStart := strings.Index(html, `<table id="team-changes-table"`)
	tableEnd := strings.Index(html, `<button class="table-expander" type="button" data-table="team-changes-table">`)
	if tableStart == -1 || tableEnd == -1 {
		t.Fatalf("expected Team changes table and preview control: %s", html)
	}
	if strings.Contains(html[tableStart:tableEnd], now.Format("2006-01-02")) {
		t.Fatalf("expected Team changes table to omit exact latest-change dates: %s", html[tableStart:tableEnd])
	}
}

func TestRenderHTMLIncludesGeneratedTimestampAndThemeControls(t *testing.T) {
	html, err := render.BuildHTML(
		map[string]models.PlayerOwnershipRecord{
			"p1": {PlayerID: "p1", PlayerName: "Haaland", ManagerCount: 5},
		},
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	timestampPattern := regexp.MustCompile(`Updated: \d{4}-\d{2}-\d{2} \d{2}:\d{2}`)
	if !timestampPattern.MatchString(html) {
		t.Fatalf("expected a labeled minute-precision generation timestamp in HTML: %s", html)
	}
	if strings.Contains(html, "Updated:") && strings.Contains(html, " UTC") {
		t.Fatalf("expected generation timestamp without a timezone label: %s", html)
	}
	if !strings.Contains(html, "<title>Fantasy Football Data</title>") || !strings.Contains(html, `<h1 class="report-title">Fantasy Football Data</h1>`) || !strings.Contains(html, `<h1 class="mobile-title report-title">Fantasy Football Data</h1>`) {
		t.Fatalf("expected Fantasy Football Data browser and report titles: %s", html)
	}
	if strings.Contains(html, "Guy Sports Data") {
		t.Fatalf("expected former Guy Sports Data title to be absent: %s", html)
	}
	if strings.Count(html, `class="report-logo" src="assets/ffd-logo.svg" alt="FFD"`) != 1 || !strings.Contains(html, `<nav class="report-nav" aria-label="Report sections"><img class="sidebar-logo" src="assets/ffd-logo.svg" alt="FFD">`) || !strings.Contains(html, `rel="icon" href="assets/ffd-favicon.svg"`) {
		t.Fatalf("expected the compact mobile logo, enlarged sidebar logo, and favicon: %s", html)
	}
	if !strings.Contains(html, `--accent:#0057b8`) || !strings.Contains(html, `.report-title{color:var(--accent)}`) || !strings.Contains(html, `.mobile-header .report-logo{width:5rem;height:3rem}`) || !strings.Contains(html, `.mobile-title{font-size:1.65rem}h2{font-size:1.35rem}`) || !strings.Contains(html, `.sidebar-logo{display:block;width:9rem;height:5.4rem`) {
		t.Fatalf("expected selected blue theme, larger mobile report title, and enlarged sidebar logo: %s", html)
	}
	if strings.Index(html, "Fantasy Football Data") > strings.Index(html, "Updated:") {
		t.Fatalf("expected generation timestamp directly below report title: %s", html)
	}
	if strings.Count(html, `class="theme-toggle"`) != 1 || !strings.Contains(html, `<div class="nav-settings"><span class="nav-label">Settings</span>`) {
		t.Fatalf("expected one theme toggle in navigation settings: %s", html)
	}
	if !strings.Contains(html, `prefers-color-scheme: dark`) {
		t.Fatalf("expected device theme preference detection in HTML: %s", html)
	}
	if !strings.Contains(html, `localStorage.getItem('report-theme')`) || !strings.Contains(html, `localStorage.setItem('report-theme',next)`) {
		t.Fatalf("expected saved theme preference support in HTML: %s", html)
	}
	if !strings.Contains(html, `:root[data-theme="dark"]`) || !strings.Contains(html, `--chart-bg`) || !strings.Contains(html, `--surface-muted`) {
		t.Fatalf("expected dark theme styles for report surfaces in HTML: %s", html)
	}
}

func TestHistoricalTrendsChartPositionsXAxisByElapsedTime(t *testing.T) {
	// Snapshots captured with highly irregular time gaps (1 hour, then 10 days, then 30 minutes)
	current := map[string]models.PlayerOwnershipRecord{
		"p1": {PlayerID: "p1", PlayerName: "Haaland", ManagerCount: 20},
		"p2": {PlayerID: "p2", PlayerName: "Salah", ManagerCount: 15},
	}
	history := map[string][]models.PlayerOwnershipHistoryEntry{
		"p1": {
			{CapturedAt: "2026-01-01T10:00:00Z", ManagerCount: 10},
			{CapturedAt: "2026-01-01T11:00:00Z", ManagerCount: 12}, // +1 hour
			{CapturedAt: "2026-01-11T11:00:00Z", ManagerCount: 18}, // +10 days
			{CapturedAt: "2026-01-11T11:30:00Z", ManagerCount: 20}, // +30 minutes
		},
		"p2": {
			{CapturedAt: "2026-01-01T10:00:00Z", ManagerCount: 8},
			{CapturedAt: "2026-01-01T11:00:00Z", ManagerCount: 10},
			{CapturedAt: "2026-01-11T11:00:00Z", ManagerCount: 14},
			{CapturedAt: "2026-01-11T11:30:00Z", ManagerCount: 15},
		},
	}

	html, err := render.BuildHTML(current, history, nil)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	// Extract the tick line x1 coordinates on the x-axis: <line x1="..." y1="700.0" x2="..." y2="706.0"
	re := regexp.MustCompile(`<line x1="([0-9.]+)" y1="700\.0" x2="[0-9.]+" y2="706\.0"`)
	matches := re.FindAllStringSubmatch(html, -1)
	if len(matches) != 4 {
		t.Fatalf("expected 4 x-axis ticks for 4 capture timestamps, got %d", len(matches))
	}

	xCoords := make([]float64, 4)
	for i, m := range matches {
		val, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			t.Fatalf("failed to parse x coordinate %q: %v", m[1], err)
		}
		xCoords[i] = val
	}

	if xCoords[0] != 60 || xCoords[3] != 900 {
		t.Fatalf("expected first and last timestamps at plot bounds, got %v", xCoords)
	}

	firstGap := xCoords[1] - xCoords[0]
	longGap := xCoords[2] - xCoords[1]
	lastGap := xCoords[3] - xCoords[2]
	if !(firstGap > 0 && longGap > firstGap*100 && lastGap > 0 && lastGap < firstGap) {
		t.Fatalf("expected x-axis gaps proportional to elapsed time, got %v", xCoords)
	}
}

func TestHistoricalTrendsChartSingleTimestamp(t *testing.T) {
	current := map[string]models.PlayerOwnershipRecord{
		"p1": {PlayerID: "p1", PlayerName: "Haaland", ManagerCount: 5},
	}
	history := map[string][]models.PlayerOwnershipHistoryEntry{
		"p1": {{CapturedAt: "2026-08-29T21:41:58Z", ManagerCount: 5}},
	}

	html, err := render.BuildHTML(current, history, nil)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	if !strings.Contains(html, "<svg class=\"trend-chart\"") {
		t.Fatalf("expected SVG line chart for single timestamp")
	}
	if !strings.Contains(html, "Haaland") {
		t.Fatalf("expected Haaland in chart output")
	}
}

func TestRenderHTMLIncludesInteractiveManagerChangeDetails(t *testing.T) {
	current := map[string]models.PlayerOwnershipRecord{
		"p1": {PlayerID: "p1", PlayerName: "Haaland", ManagerCount: 5},
	}
	history := map[string][]models.PlayerOwnershipHistoryEntry{
		"p1": {{CapturedAt: "2026-08-29T21:41:58Z", ManagerCount: 5}},
	}
	summaries := []models.ManagerChangeSummary{
		{
			ManagerID:                "m1",
			ManagerName:              "John_Doe",
			TeamName:                 "Finsbury_Bridge",
			TotalChanges:             1,
			LatestChangeAt:           "2026-08-30T18:00:00Z",
			ChangedSinceLastSnapshot: true,
			EventHistory: []models.TeamChangeEvent{
				{
					ManagerID:      "m1",
					FromCapturedAt: "2026-08-29T12:00:00Z",
					ToCapturedAt:   "2026-08-30T18:00:00Z",
					ChangeCount:    1,
					AddedPlayers:   []models.PlayerReference{{PlayerID: "p3", Name: "Saka"}},
					RemovedPlayers: []models.PlayerReference{{PlayerID: "p2", Name: "Salah"}},
				},
			},
		},
	}

	html, err := render.BuildHTML(current, history, summaries)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	// Verify interactive toggle function script exists
	if !strings.Contains(html, "function toggleManagerDetail") {
		t.Fatalf("expected toggleManagerDetail JavaScript function in HTML output")
	}

	// Verify manager detail row exists
	if !strings.Contains(html, "class=\"manager-detail-row\"") {
		t.Fatalf("expected manager-detail-row element in HTML output")
	}

	// Verify added and removed players with names are displayed
	if !strings.Contains(html, "+ Saka") {
		t.Fatalf("expected '+ Saka' in added players detail view")
	}
	if !strings.Contains(html, "- Salah") {
		t.Fatalf("expected '- Salah' in removed players detail view")
	}
	if !strings.Contains(html, "Change at 2026-08-30 19:00") {
		t.Fatalf("expected London-local event timestamp: %s", html)
	}
	if !strings.Contains(html, "Haaland: 5 managers at 2026-08-29 22:41") {
		t.Fatalf("expected London-local chart tooltip timestamp: %s", html)
	}
	for _, expected := range []string{
		`@media(min-width:1000px){body{max-width:1150px`,
		`#team-changes-table thead th:first-child,#team-changes-table tbody tr:not(.manager-detail-row) td:first-child{display:none}`,
		`.report-nav{inset:var(--mobile-header-height,4.25rem) 0 auto 0;max-height:`,
		`updateMobileHeaderHeight`,
		`header.getBoundingClientRect().height`,
	} {
		if !strings.Contains(html, expected) {
			t.Fatalf("expected mobile report layout rule %q: %s", expected, html)
		}
	}
}

func TestRenderHTMLIncludesReversibleTablePreviews(t *testing.T) {
	current := make(map[string]models.PlayerOwnershipRecord)
	for index := 0; index < 11; index++ {
		playerID := fmt.Sprintf("p%d", index)
		current[playerID] = models.PlayerOwnershipRecord{PlayerID: playerID, PlayerName: fmt.Sprintf("Player %d", index), ManagerCount: 11 - index}
	}

	summaries := make([]models.ManagerChangeSummary, 11)
	for index := range summaries {
		summaries[index] = models.ManagerChangeSummary{ManagerID: fmt.Sprintf("m%d", index), ManagerName: fmt.Sprintf("Manager %d", index), TeamName: fmt.Sprintf("Team %d", index)}
	}

	html, err := render.BuildHTML(current, nil, summaries)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	for _, expected := range []string{
		`id="player-ownership-table" class="table-preview"`,
		`id="team-changes-table" class="table-preview"`,
		`data-table="player-ownership-table"`,
		`data-table="team-changes-table"`,
		`var hiddenRows=rows.slice(10)`,
		`button.textContent=expanded?'Show fewer '+itemLabel:'Show all '+rows.length+' '+itemLabel`,
		`.manager-detail-row`,
	} {
		if !strings.Contains(html, expected) {
			t.Fatalf("expected reversible table preview markup %q in HTML: %s", expected, html)
		}
	}
}
