package unit

import (
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
}
