package processing

import (
	"fmt"
	"sort"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

// BuildOwnership produces the latest ownership counts and the historical trend by capture time.
func BuildOwnership(snapshots []models.ManagerSnapshot) (map[string]models.PlayerOwnershipRecord, map[string][]models.PlayerOwnershipHistoryEntry, error) {
	if len(snapshots) == 0 {
		return map[string]models.PlayerOwnershipRecord{}, map[string][]models.PlayerOwnershipHistoryEntry{}, nil
	}

	latestByManager := map[string]models.ManagerSnapshot{}
	for _, snap := range snapshots {
		current, ok := latestByManager[snap.ManagerID]
		if !ok || snap.CapturedAt > current.CapturedAt {
			latestByManager[snap.ManagerID] = snap
		}
	}

	current := map[string]models.PlayerOwnershipRecord{}
	for _, snap := range latestByManager {
		for _, player := range snap.Players {
			record, ok := current[player.PlayerID]
			if !ok {
				record = models.PlayerOwnershipRecord{PlayerID: player.PlayerID, PlayerName: player.Name, CapturedAt: snap.CapturedAt, ManagerCount: 0}
			}
			record.ManagerCount++
			record.CapturedAt = snap.CapturedAt
			if record.PlayerName == "" {
				record.PlayerName = player.Name
			}
			current[player.PlayerID] = record
		}
	}

	timestamps := map[string][]models.ManagerSnapshot{}
	for _, snap := range snapshots {
		timestamps[snap.CapturedAt] = append(timestamps[snap.CapturedAt], snap)
	}

	history := map[string][]models.PlayerOwnershipHistoryEntry{}
	keys := make([]string, 0, len(timestamps))
	for key := range timestamps {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		t1, err1 := time.Parse(time.RFC3339, keys[i])
		t2, err2 := time.Parse(time.RFC3339, keys[j])
		if err1 != nil || err2 != nil {
			return keys[i] < keys[j]
		}
		return t1.Before(t2)
	})

	for _, capturedAt := range keys {
		counts := map[string]int{}
		seenManagers := map[string]bool{}
		for _, snap := range timestamps[capturedAt] {
			if seenManagers[snap.ManagerID] {
				continue
			}
			seenManagers[snap.ManagerID] = true
			for _, player := range snap.Players {
				counts[player.PlayerID]++
			}
		}
		for playerID, count := range counts {
			history[playerID] = append(history[playerID], models.PlayerOwnershipHistoryEntry{CapturedAt: capturedAt, ManagerCount: count})
		}
	}

	return current, history, nil
}

// BuildChangeSummaries compares consecutive snapshots for each manager and summarizes the delta history.
func BuildChangeSummaries(snapshots []models.ManagerSnapshot) ([]models.ManagerChangeSummary, error) {
	if len(snapshots) == 0 {
		return nil, nil
	}

	byManager := map[string][]models.ManagerSnapshot{}
	for _, snap := range snapshots {
		byManager[snap.ManagerID] = append(byManager[snap.ManagerID], snap)
	}

	out := make([]models.ManagerChangeSummary, 0, len(byManager))
	for managerID, entries := range byManager {
		sort.Slice(entries, func(i, j int) bool {
			t1, err1 := time.Parse(time.RFC3339, entries[i].CapturedAt)
			t2, err2 := time.Parse(time.RFC3339, entries[j].CapturedAt)
			if err1 != nil || err2 != nil {
				return entries[i].CapturedAt < entries[j].CapturedAt
			}
			return t1.Before(t2)
		})
		summary := models.ManagerChangeSummary{ManagerID: managerID}
		var last map[string]bool
		for i := 0; i < len(entries); i++ {
			snap := entries[i]
			currentPlayers := map[string]bool{}
			for _, player := range snap.Players {
				currentPlayers[player.PlayerID] = true
			}
			if i > 0 {
				event := models.TeamChangeEvent{ManagerID: managerID, FromCapturedAt: entries[i-1].CapturedAt, ToCapturedAt: snap.CapturedAt}
				for id := range currentPlayers {
					if !last[id] {
						event.AddedPlayers = append(event.AddedPlayers, models.PlayerReference{PlayerID: id})
					}
				}
				for id := range last {
					if !currentPlayers[id] {
						event.RemovedPlayers = append(event.RemovedPlayers, models.PlayerReference{PlayerID: id})
					}
				}
				event.ChangeCount = len(event.AddedPlayers) + len(event.RemovedPlayers)
				if event.ChangeCount > 0 {
					summary.TotalChanges++
					summary.LatestChangeAt = snap.CapturedAt
					summary.ChangedSinceLastSnapshot = true
					summary.EventHistory = append(summary.EventHistory, event)
				}
			}
			last = currentPlayers
		}
		if len(entries) > 0 {
			latestSnap := entries[len(entries)-1]
			summary.ManagerName = latestSnap.ManagerName
			summary.TeamName = latestSnap.TeamName
			if summary.TeamName == "" {
				summary.TeamName = summary.ManagerName
			}
		}
		if summary.TotalChanges == 0 {
			summary.ChangedSinceLastSnapshot = false
			summary.LatestChangeAt = ""
		}
		out = append(out, summary)
	}

	return out, nil
}

func init() {
	_ = fmt.Sprintf
}
