package unit

import (
	"testing"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/processing"
)

func TestBuildManagerChangeSummaries(t *testing.T) {
	base := time.Date(2026, 1, 10, 9, 0, 0, 0, time.UTC)
	snapshots := []models.ManagerSnapshot{
		{ManagerID: "m1", ManagerName: "A", CapturedAt: base.Add(0).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p2"}}},
		{ManagerID: "m1", ManagerName: "A", CapturedAt: base.Add(3 * time.Hour).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p3"}}},
		{ManagerID: "m1", ManagerName: "A", CapturedAt: base.Add(6 * time.Hour).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p3"}}},
	}

	summaries, err := processing.BuildChangeSummaries(snapshots)
	if err != nil {
		t.Fatalf("BuildChangeSummaries returned error: %v", err)
	}

	if len(summaries) != 1 {
		t.Fatalf("expected one summary, got %d", len(summaries))
	}
	if summaries[0].TotalChanges != 1 {
		t.Fatalf("expected total changes 1, got %d", summaries[0].TotalChanges)
	}
	if !summaries[0].ChangedSinceLastSnapshot {
		t.Fatal("expected manager to be marked changed since last snapshot")
	}
}
