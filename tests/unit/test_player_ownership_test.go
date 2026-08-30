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

func TestPartialCaptureRunPreservesUnchangedManagerCounts(t *testing.T) {
	// Scenario: 28 managers at T1, all with Haaland (p1).
	// At T2, only 2 managers (m1, m2) have new snapshots. The other 26 managers have NO new snapshot at T2.
	t1 := "2026-08-29T12:00:00Z"
	t2 := "2026-08-30T18:00:00Z"

	snapshots := make([]models.ManagerSnapshot, 0)
	// Create T1 snapshots for 28 managers
	for i := 1; i <= 28; i++ {
		mID := string(rune('a'+(i%26))) + string(rune('0'+(i/26)))
		snapshots = append(snapshots, models.ManagerSnapshot{
			ManagerID:   mID,
			ManagerName: "Mgr_" + mID,
			CapturedAt:  t1,
			Players:     []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}},
		})
	}

	// At T2, only manager m1 and m2 have new snapshots
	snapshots = append(snapshots, models.ManagerSnapshot{
		ManagerID:   snapshots[0].ManagerID,
		ManagerName: snapshots[0].ManagerName,
		CapturedAt:  t2,
		Players:     []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}, {PlayerID: "p2", Name: "Salah"}},
	})
	snapshots = append(snapshots, models.ManagerSnapshot{
		ManagerID:   snapshots[1].ManagerID,
		ManagerName: snapshots[1].ManagerName,
		CapturedAt:  t2,
		Players:     []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}, {PlayerID: "p3", Name: "Saka"}},
	})

	current, history, err := processing.BuildOwnership(snapshots)
	if err != nil {
		t.Fatalf("BuildOwnership failed: %v", err)
	}

	// Haaland (p1) must have count 28 at T1 AND count 28 at T2 (NOT 2!)
	if got := current["p1"].ManagerCount; got != 28 {
		t.Fatalf("expected current Haaland count 28, got %d", got)
	}

	p1History := history["p1"]
	if len(p1History) != 2 {
		t.Fatalf("expected 2 history entries for p1, got %d", len(p1History))
	}

	if p1History[0].ManagerCount != 28 {
		t.Fatalf("expected Haaland count 28 at T1, got %d", p1History[0].ManagerCount)
	}

	if p1History[1].ManagerCount != 28 {
		t.Fatalf("expected Haaland count 28 at T2 (partial capture run), got %d", p1History[1].ManagerCount)
	}
}
