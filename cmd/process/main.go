package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/config"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/processing"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
	"gopkg.in/yaml.v3"
)

func main() {
	cfg := config.Load()
	inputDir := cfg.RawDir
	changesOutDir := filepath.Join(cfg.ProcessDir, "manager-changes")

	if err := os.MkdirAll(cfg.ProcessDir, 0o755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(changesOutDir, 0o755); err != nil {
		panic(err)
	}

	snapshots, err := storage.SnapshotFiles(inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process failed: %v\n", err)
		os.Exit(1)
	}

	current, history, err := processing.BuildOwnership(snapshots)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ownership build failed: %v\n", err)
		os.Exit(1)
	}

	// Build unified player ownership list
	playerList := make([]models.PlayerOwnershipDetail, 0, len(current))
	for id, rec := range current {
		playerList = append(playerList, models.PlayerOwnershipDetail{
			PlayerID:     id,
			PlayerName:   rec.PlayerName,
			CurrentCount: rec.ManagerCount,
			History:      history[id],
		})
	}

	// Sort player ownership list by CurrentCount descending, then PlayerName ascending
	sort.Slice(playerList, func(i, j int) bool {
		if playerList[i].CurrentCount != playerList[j].CurrentCount {
			return playerList[i].CurrentCount > playerList[j].CurrentCount
		}
		return playerList[i].PlayerName < playerList[j].PlayerName
	})

	doc := models.PlayerOwnershipDoc{
		Season:      cfg.Season,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Players:     playerList,
	}

	docData, err := yaml.Marshal(doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal player ownership failed: %v\n", err)
		os.Exit(1)
	}

	ownershipFilePath := filepath.Join(cfg.ProcessDir, "player-ownership.yaml")
	if err := os.WriteFile(ownershipFilePath, docData, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write player-ownership.yaml failed: %v\n", err)
		os.Exit(1)
	}

	// Write manager changes
	changes, err := processing.BuildChangeSummaries(snapshots)
	if err != nil {
		fmt.Fprintf(os.Stderr, "change summary failed: %v\n", err)
		os.Exit(1)
	}

	for _, summary := range changes {
		data, err := yaml.Marshal(summary)
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal manager change failed for %s: %v\n", summary.ManagerID, err)
			continue
		}
		teamName := summary.TeamName
		if teamName == "" {
			teamName = summary.ManagerName
		}
		if teamName == "" {
			teamName = fmt.Sprintf("team_%s", summary.ManagerID)
		}
		filename := fmt.Sprintf("%s_%s.yaml", sanitizeFilename(teamName), sanitizeFilename(summary.ManagerID))
		if err := os.WriteFile(filepath.Join(changesOutDir, filename), data, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write changes failed: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("processed %d snapshot records into %s and %s\n", len(snapshots), ownershipFilePath, changesOutDir)
}

func sanitizeFilename(s string) string {
	return strings.NewReplacer(" ", "_", "/", "_", "?", "_", "&", "_", "=", "_", "\\", "_", ":", "_").Replace(s)
}
