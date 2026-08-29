package unit

import (
	"strings"
	"testing"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/render"
)

func TestRenderHTMLIncludesOwnershipAndChanges(t *testing.T) {
	current := map[string]models.PlayerOwnershipRecord{
		"p1": {PlayerID: "p1", PlayerName: "Haaland", ManagerCount: 2},
		"p2": {PlayerID: "p2", PlayerName: "Salah", ManagerCount: 1},
	}
	history := map[string][]models.PlayerOwnershipHistoryEntry{
		"p1": {{CapturedAt: "2026-01-10T12:00:00Z", ManagerCount: 1}, {CapturedAt: "2026-01-10T14:00:00Z", ManagerCount: 2}},
	}
	summaries := []models.ManagerChangeSummary{{ManagerID: "m1", TotalChanges: 1, LatestChangeAt: "2026-01-10T15:00:00Z", ChangedSinceLastSnapshot: true}}

	html, err := render.BuildHTML(current, history, summaries)
	if err != nil {
		t.Fatalf("BuildHTML returned error: %v", err)
	}

	if !strings.Contains(html, "Player ownership") || !strings.Contains(html, "Haaland") || !strings.Contains(html, "m1") {
		t.Fatalf("rendered output missing expected content: %s", html)
	}
}
