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

	// Group snapshots by timestamp
	timestamps := map[string][]models.ManagerSnapshot{}
	for _, snap := range snapshots {
		timestamps[snap.CapturedAt] = append(timestamps[snap.CapturedAt], snap)
	}

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

	// Maintain active manager state as of each timestamp T.
	// Because unchanged manager snapshots are skipped during capture,
	// a manager's active team at timestamp T is their latest snapshot captured on or before T.
	activeSnapshots := map[string]models.ManagerSnapshot{}
	playerNames := map[string]string{}
	history := map[string][]models.PlayerOwnershipHistoryEntry{}

	for _, capturedAt := range keys {
		// Update active snapshot for each manager who captured at this timestamp T
		for _, snap := range timestamps[capturedAt] {
			activeSnapshots[snap.ManagerID] = snap
			for _, player := range snap.Players {
				if player.Name != "" {
					playerNames[player.PlayerID] = player.Name
				}
			}
		}

		// Count player frequencies across ALL active manager snapshots in the league as of timestamp T
		counts := map[string]int{}
		for _, snap := range activeSnapshots {
			for _, player := range snap.Players {
				counts[player.PlayerID]++
			}
		}

		// Record the point-in-time count for each player present at timestamp T
		for playerID, count := range counts {
			history[playerID] = append(history[playerID], models.PlayerOwnershipHistoryEntry{
				CapturedAt:   capturedAt,
				ManagerCount: count,
			})
		}
	}

	// Current ownership is derived from activeSnapshots at the final (latest) timestamp
	current := map[string]models.PlayerOwnershipRecord{}
	counts := map[string]int{}
	latestTimestamp := keys[len(keys)-1]

	for _, snap := range activeSnapshots {
		for _, player := range snap.Players {
			counts[player.PlayerID]++
			if player.Name != "" {
				playerNames[player.PlayerID] = player.Name
			}
		}
	}

	for playerID, count := range counts {
		current[playerID] = models.PlayerOwnershipRecord{
			PlayerID:     playerID,
			PlayerName:   playerNames[playerID],
			CapturedAt:   latestTimestamp,
			ManagerCount: count,
		}
	}

	return current, history, nil
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// BuildChangeSummaries compares consecutive snapshots for each manager and summarizes the delta history.
func BuildChangeSummaries(snapshots []models.ManagerSnapshot) ([]models.ManagerChangeSummary, error) {
	if len(snapshots) == 0 {
		return nil, nil
	}

	byManager := map[string][]models.ManagerSnapshot{}
	playerNames := map[string]string{}

	for _, snap := range snapshots {
		byManager[snap.ManagerID] = append(byManager[snap.ManagerID], snap)
		for _, p := range snap.Players {
			if p.Name != "" {
				playerNames[p.PlayerID] = p.Name
			}
		}
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
				if player.Name != "" {
					playerNames[player.PlayerID] = player.Name
				}
			}
			if i > 0 {
				event := models.TeamChangeEvent{
					ManagerID:      managerID,
					FromCapturedAt: entries[i-1].CapturedAt,
					ToCapturedAt:   snap.CapturedAt,
				}
				for id := range currentPlayers {
					if !last[id] {
						name := playerNames[id]
						event.AddedPlayers = append(event.AddedPlayers, models.PlayerReference{
							PlayerID: id,
							Name:     name,
						})
					}
				}
				for id := range last {
					if !currentPlayers[id] {
						name := playerNames[id]
						event.RemovedPlayers = append(event.RemovedPlayers, models.PlayerReference{
							PlayerID: id,
							Name:     name,
						})
					}
				}

				// Sort added and removed players by Name/PlayerID for deterministic output
				sort.Slice(event.AddedPlayers, func(a, b int) bool {
					return event.AddedPlayers[a].PlayerID < event.AddedPlayers[b].PlayerID
				})
				sort.Slice(event.RemovedPlayers, func(a, b int) bool {
					return event.RemovedPlayers[a].PlayerID < event.RemovedPlayers[b].PlayerID
				})

				// A substitution/transfer (1 player removed, 1 player added) equals 1 change
				event.ChangeCount = max(len(event.AddedPlayers), len(event.RemovedPlayers))

				if event.ChangeCount > 0 {
					summary.TotalChanges += event.ChangeCount
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
