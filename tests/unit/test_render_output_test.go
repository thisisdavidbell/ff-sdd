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

func TestHistoricalTrendsChartHasConsistentXAxisGaps(t *testing.T) {
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

	// Verify the gap between consecutive ticks is constant across all 4 dates
	expectedGap := xCoords[1] - xCoords[0]
	if expectedGap <= 0 {
		t.Fatalf("expected positive gap between ticks, got %f", expectedGap)
	}

	for i := 1; i < len(xCoords)-1; i++ {
		gap := xCoords[i+1] - xCoords[i]
		diff := gap - expectedGap
		if diff < -0.01 || diff > 0.01 {
			t.Fatalf("inconsistent gap between tick %d and %d: got %f, expected %f (all xCoords: %v)", i, i+1, gap, expectedGap, xCoords)
		}
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
