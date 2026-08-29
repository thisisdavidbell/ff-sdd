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
	DetailURL string
	Players   []models.PlayerReference
}

// ExtractManagersFromHTML parses the season table for publicteamdetail links.
func ExtractManagersFromHTML(page string) []ManagerSummary {
	pattern := regexp.MustCompile(`(?i)publicteamdetail\.php\?[^'"<>[:space:]]+`)
	matches := pattern.FindAllString(page, -1)
	if len(matches) == 0 {
		return nil
	}

	out := make([]ManagerSummary, 0, len(matches))
	seenIDs := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		parsed, err := url.Parse(html.UnescapeString(match))
		if err != nil {
			continue
		}
		id := strings.TrimSpace(parsed.Query().Get("id"))
		name := strings.TrimSpace(parsed.Query().Get("name"))
		if id == "" || name == "" {
			continue
		}
		if _, seen := seenIDs[id]; seen {
			continue
		}
		seenIDs[id] = struct{}{}
		detailURL := fmt.Sprintf("https://www.guysports.co.uk/guysports/publicteamdetail.php?id=%s&name=%s&view=public", id, url.QueryEscape(name))
		out = append(out, ManagerSummary{ID: id, Name: name, DetailURL: detailURL})
	}
	return out
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
