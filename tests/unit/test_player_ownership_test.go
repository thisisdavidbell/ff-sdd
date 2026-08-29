package unit

import (
	"testing"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/processing"
)

func TestBuildCurrentAndHistoricalOwnership(t *testing.T) {
	base := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	snapshots := []models.ManagerSnapshot{
		{ManagerID: "m1", ManagerName: "A", CapturedAt: base.Add(0).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}, {PlayerID: "p2", Name: "Salah"}}},
		{ManagerID: "m2", ManagerName: "B", CapturedAt: base.Add(2 * time.Hour).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}, {PlayerID: "p3", Name: "De Bruyne"}}},
		{ManagerID: "m1", ManagerName: "A", CapturedAt: base.Add(4 * time.Hour).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}, {PlayerID: "p4", Name: "Foden"}}},
	}

	current, history, err := processing.BuildOwnership(snapshots)
	if err != nil {
		t.Fatalf("BuildOwnership returned error: %v", err)
	}

	if got := current["p1"].ManagerCount; got != 2 {
		t.Fatalf("expected current p1 count 2, got %d", got)
	}
	if got := current["p2"].ManagerCount; got != 0 {
		t.Fatalf("expected p2 to be absent from latest snapshot, got %d", got)
	}
	if len(history["p1"]) != 3 {
		t.Fatalf("expected p1 historical entries to include all capture events, got %d", len(history["p1"]))
	}
}
