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
		{ManagerID: "m1", ManagerName: "John_Doe", TeamName: "Finsbury_Bridge", CapturedAt: base.Add(0).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p2"}}},
		{ManagerID: "m1", ManagerName: "John_Doe", TeamName: "Finsbury_Bridge", CapturedAt: base.Add(3 * time.Hour).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p3"}}},
		{ManagerID: "m1", ManagerName: "John_Doe", TeamName: "Finsbury_Bridge", CapturedAt: base.Add(6 * time.Hour).Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p3"}}},
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
	if summaries[0].ManagerName != "John_Doe" {
		t.Fatalf("expected manager name John_Doe, got %q", summaries[0].ManagerName)
	}
	if summaries[0].TeamName != "Finsbury_Bridge" {
		t.Fatalf("expected team name Finsbury_Bridge, got %q", summaries[0].TeamName)
	}
	if !summaries[0].ChangedSinceLastSnapshot {
		t.Fatal("expected manager to be marked changed since last snapshot")
	}
}

func TestSingleSubstitutionCountsAsOneChange(t *testing.T) {
	base := time.Date(2026, 1, 10, 9, 0, 0, 0, time.UTC)
	snapshots := []models.ManagerSnapshot{
		{
			ManagerID:   "m1",
			ManagerName: "John_Doe",
			CapturedAt:  base.Add(0).Format(time.RFC3339),
			Players:     []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}, {PlayerID: "p2", Name: "Salah"}},
		},
		{
			ManagerID:   "m1",
			ManagerName: "John_Doe",
			CapturedAt:  base.Add(24 * time.Hour).Format(time.RFC3339),
			Players:     []models.PlayerReference{{PlayerID: "p1", Name: "Haaland"}, {PlayerID: "p3", Name: "Saka"}}, // Swapped Salah for Saka
		},
	}

	summaries, err := processing.BuildChangeSummaries(snapshots)
	if err != nil {
		t.Fatalf("BuildChangeSummaries returned error: %v", err)
	}

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}

	sum := summaries[0]
	if sum.TotalChanges != 1 {
		t.Fatalf("expected TotalChanges == 1 for single substitution (1 in, 1 out), got %d", sum.TotalChanges)
	}

	if len(sum.EventHistory) != 1 {
		t.Fatalf("expected 1 event, got %d", len(sum.EventHistory))
	}

	event := sum.EventHistory[0]
	if event.ChangeCount != 1 {
		t.Fatalf("expected Event.ChangeCount == 1, got %d", event.ChangeCount)
	}

	if len(event.AddedPlayers) != 1 || event.AddedPlayers[0].Name != "Saka" {
		t.Fatalf("expected AddedPlayers to contain Saka with name, got %#v", event.AddedPlayers)
	}

	if len(event.RemovedPlayers) != 1 || event.RemovedPlayers[0].Name != "Salah" {
		t.Fatalf("expected RemovedPlayers to contain Salah with name, got %#v", event.RemovedPlayers)
	}
}
