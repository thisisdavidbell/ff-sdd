package validation

import (
	"fmt"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

// ValidateSnapshot enforces the minimum contract for a historical manager snapshot.
func ValidateSnapshot(snapshot models.ManagerSnapshot) error {
	if snapshot.ManagerID == "" {
		return fmt.Errorf("manager_id is required")
	}
	if snapshot.CapturedAt == "" {
		return fmt.Errorf("captured_at is required")
	}
	if _, err := time.Parse(time.RFC3339, snapshot.CapturedAt); err != nil {
		return fmt.Errorf("captured_at must be RFC3339: %w", err)
	}
	if snapshot.Source == "" {
		return fmt.Errorf("source is required")
	}
	seen := map[string]bool{}
	for _, player := range snapshot.Players {
		if player.PlayerID == "" {
			return fmt.Errorf("player_id is required for snapshot %s", snapshot.ManagerID)
		}
		if seen[player.PlayerID] {
			return fmt.Errorf("duplicate player_id %s in snapshot %s", player.PlayerID, snapshot.ManagerID)
		}
		seen[player.PlayerID] = true
	}
	return nil
}
