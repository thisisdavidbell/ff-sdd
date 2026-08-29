package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/capture"
	"github.com/thisisdavidbell/ff-sdd/internal/config"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
)

func main() {
	cfg := config.Load()
	source := flag.String("source", "guysports", "source name")
	output := flag.String("output", cfg.RawDir, "output directory")
	pagesFlag := flag.Int("pages", cfg.Pages, "number of pages to capture")
	flag.Parse()

	if _, err := os.Stat(*output); err != nil {
		if err := os.MkdirAll(*output, 0o755); err != nil {
			panic(err)
		}
	}

	if *source != "guysports" {
		fmt.Fprintf(os.Stderr, "unsupported source: %s\n", *source)
		os.Exit(2)
	}

	totalPages := *pagesFlag
	if totalPages <= 0 {
		totalPages = 1
	}

	allManagers := make([]capture.ManagerSummary, 0)
	seenIDs := make(map[string]bool)

	baseURL := cfg.SourceURL
	if !strings.Contains(baseURL, "?") {
		baseURL = strings.TrimRight(baseURL, "/")
	}

	for p := 1; p <= totalPages; p++ {
		pageURL := fmt.Sprintf("%s?page=%d", strings.Split(baseURL, "?")[0], p)
		resp, err := http.Get(pageURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch page %d failed: %v\n", p, err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "read body page %d failed: %v\n", p, err)
			continue
		}
		html := string(body)
		if !strings.Contains(html, "GuySports") && !strings.Contains(html, "Season Table") {
			fmt.Fprintf(os.Stderr, "unexpected page payload for page %d\n", p)
			continue
		}

		pageManagers := capture.ExtractManagersFromHTML(html)
		for _, m := range pageManagers {
			if !seenIDs[m.ID] {
				seenIDs[m.ID] = true
				allManagers = append(allManagers, m)
			}
		}
	}

	if len(allManagers) == 0 {
		fmt.Fprintf(os.Stderr, "no managers found across %d pages\n", totalPages)
		os.Exit(1)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	capturedCount := 0
	skippedCount := 0

	for _, manager := range allManagers {
		players, err := capture.FetchAndParseManagerDetail(manager.DetailURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch manager detail failed for %s: %v\n", manager.Name, err)
			continue
		}

		teamName := capture.FormatWithUnderscores(manager.TeamName)
		if teamName == "" {
			teamName = capture.FormatWithUnderscores(manager.Name)
		}
		if teamName == "" {
			teamName = fmt.Sprintf("Team_%s", manager.ID)
		}

		managerName := capture.FormatWithUnderscores(manager.Manager)
		if managerName == "" {
			managerName = teamName
		}

		managerDir := filepath.Join(*output, fmt.Sprintf("%s_%s", sanitizeFilename(teamName), sanitizeFilename(manager.ID)))
		if err := os.MkdirAll(managerDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir manager dir failed: %v\n", err)
			continue
		}

		// Check if latest snapshot is unchanged
		if isLatestSnapshotUnchanged(managerDir, players) {
			skippedCount++
			continue
		}

		snapshot := models.ManagerSnapshot{
			ManagerID:   manager.ID,
			ManagerName: managerName,
			TeamName:    teamName,
			CapturedAt:  now,
			Source:      *source,
			Players:     players,
			RawURL:      manager.DetailURL,
		}

		path := filepath.Join(managerDir, fmt.Sprintf("%s.yaml", strings.ReplaceAll(now, ":", "-")))
		if err := capture.WriteSnapshot(path, snapshot); err != nil {
			fmt.Fprintf(os.Stderr, "write snapshot failed: %v\n", err)
			os.Exit(1)
		}
		capturedCount++
	}

	fmt.Printf("captured %d managers (%d unchanged skipped) across %d pages into %s\n", capturedCount, skippedCount, totalPages, *output)
}

func isLatestSnapshotUnchanged(managerDir string, currentPlayers []models.PlayerReference) bool {
	entries, err := os.ReadDir(managerDir)
	if err != nil || len(entries) == 0 {
		return false
	}
	yamlFiles := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml")) {
			yamlFiles = append(yamlFiles, entry.Name())
		}
	}
	if len(yamlFiles) == 0 {
		return false
	}
	sort.Strings(yamlFiles)
	latestFile := filepath.Join(managerDir, yamlFiles[len(yamlFiles)-1])
	latestSnap, err := storage.ReadSnapshotFile(latestFile)
	if err != nil {
		return false
	}
	return capture.IsRosterUnchanged(latestSnap.Players, currentPlayers)
}

func sanitizeFilename(s string) string {
	return strings.NewReplacer(" ", "_", "/", "_", "?", "_", "&", "_", "=", "_", "\\", "_", ":", "_").Replace(s)
}
