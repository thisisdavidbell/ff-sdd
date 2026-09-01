package unit

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

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
	if !strings.Contains(html, "viewBox=\"0 0 1100 800\"") {
		t.Fatalf("expected 800px tall chart viewBox in HTML: %s", html)
	}
	if !strings.Contains(html, "class=\"player-line\"") {
		t.Fatalf("expected player line paths in HTML: %s", html)
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

	timestampPattern := regexp.MustCompile(`Report generated: \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} UTC`)
	if !timestampPattern.MatchString(html) {
		t.Fatalf("expected a labeled UTC generation timestamp in HTML: %s", html)
	}
	if strings.Index(html, "Guy Sports Team Report") > strings.Index(html, "Report generated:") {
		t.Fatalf("expected generation timestamp directly below report title: %s", html)
	}
	if !strings.Contains(html, `class="theme-toggle"`) {
		t.Fatalf("expected visible theme toggle in HTML: %s", html)
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
}
