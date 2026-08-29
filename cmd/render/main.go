package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thisisdavidbell/ff-sdd/internal/config"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/render"
	"gopkg.in/yaml.v3"
)

func main() {
	cfg := config.Load()
	ownershipFile := filepath.Join(cfg.ProcessDir, "player-ownership.yaml")
	changesDir := filepath.Join(cfg.ProcessDir, "manager-changes")
	outputFile := cfg.RenderFile

	current := map[string]models.PlayerOwnershipRecord{}
	history := map[string][]models.PlayerOwnershipHistoryEntry{}

	// Read consolidated player-ownership.yaml first
	if data, err := os.ReadFile(ownershipFile); err == nil {
		var doc models.PlayerOwnershipDoc
		if err := yaml.Unmarshal(data, &doc); err == nil && len(doc.Players) > 0 {
			for _, p := range doc.Players {
				current[p.PlayerID] = models.PlayerOwnershipRecord{
					PlayerID:     p.PlayerID,
					PlayerName:   p.PlayerName,
					ManagerCount: p.CurrentCount,
				}
				history[p.PlayerID] = p.History
			}
		}
	}

	// Fallback to reading legacy directory if single file was empty or missing
	if len(current) == 0 {
		legacyInputDir := filepath.Join(cfg.ProcessDir, "player-ownership")
		files, err := os.ReadDir(legacyInputDir)
		if err == nil {
			for _, file := range files {
				if file.IsDir() || filepath.Ext(file.Name()) != ".yaml" {
					continue
				}
				path := filepath.Join(legacyInputDir, file.Name())
				var rec models.PlayerOwnershipRecord
				data, _ := os.ReadFile(path)
				_ = yaml.Unmarshal(data, &rec)
				if rec.PlayerID != "" {
					current[rec.PlayerID] = rec
				}
			}
		}
	}

	if len(current) == 0 {
		fmt.Fprintf(os.Stderr, "no processed player ownership data found in %s; please run process first\n", cfg.ProcessDir)
		os.Exit(1)
	}

	changes := []models.ManagerChangeSummary{}
	changeFiles, err := os.ReadDir(changesDir)
	if err == nil {
		for _, file := range changeFiles {
			if file.IsDir() || (filepath.Ext(file.Name()) != ".yaml" && filepath.Ext(file.Name()) != ".yml") {
				continue
			}
			path := filepath.Join(changesDir, file.Name())
			var summary models.ManagerChangeSummary
			data, _ := os.ReadFile(path)
			_ = yaml.Unmarshal(data, &summary)
			if summary.ManagerID != "" {
				changes = append(changes, summary)
			}
		}
	}

	html, err := render.BuildHTML(current, history, changes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render failed: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "unable to prepare output dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputFile, []byte(html), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "unable to write render output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("rendered report to %s\n", outputFile)
}
