package models

// PlayerReference holds one footballer referenced in a manager's team snapshot.
type PlayerReference struct {
	PlayerID string `yaml:"player_id"`
	Name     string `yaml:"name,omitempty"`
	Position string `yaml:"position,omitempty"`
	TeamName string `yaml:"team_name,omitempty"`
}

// ManagerSnapshot records the full roster for one manager at one capture moment.
type ManagerSnapshot struct {
	ManagerID   string           `yaml:"manager_id"`
	ManagerName string           `yaml:"manager_name,omitempty"`
	CapturedAt  string           `yaml:"captured_at"`
	Source      string           `yaml:"source,omitempty"`
	Players     []PlayerReference `yaml:"players"`
	RawURL      string           `yaml:"raw_url,omitempty"`
}

// TeamChangeEvent records the delta between two consecutive snapshots for the same manager.
type TeamChangeEvent struct {
	ManagerID       string           `yaml:"manager_id"`
	FromCapturedAt  string           `yaml:"from_captured_at"`
	ToCapturedAt    string           `yaml:"to_captured_at"`
	AddedPlayers    []PlayerReference `yaml:"added_players"`
	RemovedPlayers  []PlayerReference `yaml:"removed_players"`
	ChangeCount     int              `yaml:"change_count"`
}

// PlayerOwnershipRecord is the current-count view for a player in the latest snapshot set.
type PlayerOwnershipRecord struct {
	PlayerID    string `yaml:"player_id"`
	PlayerName  string `yaml:"player_name,omitempty"`
	CapturedAt  string `yaml:"captured_at,omitempty"`
	ManagerCount int   `yaml:"manager_count"`
}

// PlayerOwnershipHistoryEntry tracks a player's count over time.
type PlayerOwnershipHistoryEntry struct {
	CapturedAt   string `yaml:"captured_at"`
	ManagerCount int    `yaml:"manager_count"`
}

// ManagerChangeSummary is the season-level summary for a manager.
type ManagerChangeSummary struct {
	ManagerID              string         `yaml:"manager_id"`
	TotalChanges           int            `yaml:"total_changes"`
	LatestChangeAt         string         `yaml:"latest_change_at,omitempty"`
	ChangedSinceLastSnapshot bool          `yaml:"changed_since_last_snapshot"`
	EventHistory           []TeamChangeEvent `yaml:"event_history,omitempty"`
}
