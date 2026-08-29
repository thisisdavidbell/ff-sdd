package capture

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
)

// ManagerSummary is a live manager reference from the season page.
type ManagerSummary struct {
	ID        string
	Name      string
	TeamName  string
	Manager   string
	DetailURL string
	Players   []models.PlayerReference
}

// FormatWithUnderscores replaces whitespace sequences with a single underscore.
func FormatWithUnderscores(s string) string {
	s = strings.TrimSpace(s)
	fields := strings.Fields(s)
	return strings.Join(fields, "_")
}

// ExtractManagersFromHTML parses the season table for publicteamdetail links.
func ExtractManagersFromHTML(page string) []ManagerSummary {
	// First attempt row-based extraction to capture both team and manager names if present in table columns
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>(.*?)</tr>`)
	rows := rowPattern.FindAllString(page, -1)

	out := make([]ManagerSummary, 0)
	seenIDs := make(map[string]struct{})

	for _, row := range rows {
		if !strings.Contains(strings.ToLower(row), "publicteamdetail.php") {
			continue
		}
		linkPattern := regexp.MustCompile(`(?i)publicteamdetail\.php\?[^'"<>[:space:]]+`)
		linkMatch := linkPattern.FindString(row)
		if linkMatch == "" {
			continue
		}
		parsed, err := url.Parse(html.UnescapeString(linkMatch))
		if err != nil {
			continue
		}
		id := strings.TrimSpace(parsed.Query().Get("id"))
		urlName := strings.TrimSpace(parsed.Query().Get("name"))
		if id == "" {
			continue
		}
		if _, seen := seenIDs[id]; seen {
			continue
		}
		seenIDs[id] = struct{}{}

		// Extract cell texts if multiple columns exist
		cellPattern := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
		cells := cellPattern.FindAllStringSubmatch(row, -1)

		teamName := urlName
		managerName := urlName

		if len(cells) >= 3 {
			c1 := stripTags(cells[1][1])
			c2 := stripTags(cells[2][1])
			if c1 != "" && !strings.EqualFold(c1, "team") && !strings.EqualFold(c1, "manager") {
				teamName = c1
			}
			if c2 != "" && !strings.EqualFold(c2, "manager") && !strings.EqualFold(c2, "team") && !strings.EqualFold(c2, "pts") && !strings.EqualFold(c2, "points") {
				managerName = c2
			}
		}

		teamFormatted := FormatWithUnderscores(teamName)
		managerFormatted := FormatWithUnderscores(managerName)
		if teamFormatted == "" {
			teamFormatted = fmt.Sprintf("Team_%s", id)
		}
		if managerFormatted == "" {
			managerFormatted = teamFormatted
		}

		detailURL := fmt.Sprintf("https://www.guysports.co.uk/guysports/publicteamdetail.php?id=%s&name=%s&view=public", id, url.QueryEscape(teamName))
		out = append(out, ManagerSummary{
			ID:        id,
			Name:      teamName,
			TeamName:  teamFormatted,
			Manager:   managerFormatted,
			DetailURL: detailURL,
		})
	}

	// Fallback to direct link extraction if no rows were matched
	if len(out) == 0 {
		pattern := regexp.MustCompile(`(?i)publicteamdetail\.php\?[^'"<>[:space:]]+`)
		matches := pattern.FindAllString(page, -1)
		for _, match := range matches {
			parsed, err := url.Parse(html.UnescapeString(match))
			if err != nil {
				continue
			}
			id := strings.TrimSpace(parsed.Query().Get("id"))
			name := strings.TrimSpace(parsed.Query().Get("name"))
			if id == "" {
				continue
			}
			if _, seen := seenIDs[id]; seen {
				continue
			}
			seenIDs[id] = struct{}{}
			formatted := FormatWithUnderscores(name)
			if formatted == "" {
				formatted = fmt.Sprintf("Team_%s", id)
			}
			detailURL := fmt.Sprintf("https://www.guysports.co.uk/guysports/publicteamdetail.php?id=%s&name=%s&view=public", id, url.QueryEscape(name))
			out = append(out, ManagerSummary{
				ID:        id,
				Name:      name,
				TeamName:  formatted,
				Manager:   formatted,
				DetailURL: detailURL,
			})
		}
	}

	return out
}

// IsRosterUnchanged checks if two player reference slices contain the exact same player IDs.
func IsRosterUnchanged(existing, current []models.PlayerReference) bool {
	if len(existing) != len(current) {
		return false
	}
	existingMap := make(map[string]bool, len(existing))
	for _, p := range existing {
		existingMap[p.PlayerID] = true
	}
	for _, p := range current {
		if !existingMap[p.PlayerID] {
			return false
		}
	}
	return true
}

// FetchAndParseManagerDetail loads the live public team page and extracts the players in the manager's team.
func FetchAndParseManagerDetail(detailURL string) ([]models.PlayerReference, error) {
	resp, err := http.Get(detailURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %s for %s", resp.Status, detailURL)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ParsePlayerRows(string(body)), nil
}

// ParsePlayerRows extracts positions and players from a live team detail page.
func ParsePlayerRows(page string) []models.PlayerReference {
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>(.*?)</tr>`)
	rows := rowPattern.FindAllString(page, -1)
	players := make([]models.PlayerReference, 0)
	for _, row := range rows {
		cellPattern := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
		cells := cellPattern.FindAllStringSubmatch(row, -1)
		if len(cells) < 3 {
			continue
		}
		position := stripTags(cells[0][1])
		name := stripTags(cells[1][1])
		teamName := stripTags(cells[2][1])
		if position == "" || name == "" || teamName == "" {
			continue
		}
		if strings.EqualFold(position, "Position") || strings.EqualFold(name, "Player") || strings.EqualFold(teamName, "Team") {
			continue
		}
		players = append(players, models.PlayerReference{
			PlayerID: sanitizePlayerID(name),
			Name:     strings.TrimSpace(name),
			Position: strings.TrimSpace(position),
			TeamName: strings.TrimSpace(teamName),
		})
	}
	return players
}

func stripTags(s string) string {
	s = regexp.MustCompile(`(?is)<[^>]+>`).ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, "\u00A0", " ")
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func sanitizePlayerID(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.NewReplacer(" ", "-", "/", "-", "&", "-", "'", "-").Replace(s)
	return s
}

// WriteSnapshot writes a manager snapshot to disk using the canonical project path layout.
func WriteSnapshot(path string, snapshot models.ManagerSnapshot) error {
	return storage.WriteSnapshotFile(path, snapshot)
}
