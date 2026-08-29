package render

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

// BuildHTML renders a static HTML summary from current ownership and change summaries.
func BuildHTML(current map[string]models.PlayerOwnershipRecord, history map[string][]models.PlayerOwnershipHistoryEntry, summaries []models.ManagerChangeSummary) (string, error) {
	// Sort players by manager count descending, then player name ascending
	playerList := make([]models.PlayerOwnershipRecord, 0, len(current))
	for _, rec := range current {
		playerList = append(playerList, rec)
	}
	sort.Slice(playerList, func(i, j int) bool {
		if playerList[i].ManagerCount != playerList[j].ManagerCount {
			return playerList[i].ManagerCount > playerList[j].ManagerCount
		}
		return playerList[i].PlayerName < playerList[j].PlayerName
	})

	var b bytes.Buffer
	b.WriteString("<!doctype html>\n")
	b.WriteString("<html><head><meta charset=\"utf-8\"><title>Guy Sports Team Report</title>\n")
	b.WriteString("<style>body{font-family:sans-serif;max-width:1100px;margin:2rem auto;padding:0 1rem}table{border-collapse:collapse;width:100%;margin:1rem 0}th,td{border:1px solid #ccc;padding:.5rem;text-align:left}th{background:#f5f5f5}</style></head><body>\n")
	b.WriteString("<h1>Guy Sports Team Report</h1>\n")
	b.WriteString("<h2>Player ownership</h2>\n")
	b.WriteString("<table><thead><tr><th>Player</th><th>Managers</th></tr></thead><tbody>\n")
	for _, rec := range playerList {
		name := rec.PlayerName
		if name == "" {
			name = rec.PlayerID
		}
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>\n", name, rec.ManagerCount))
	}
	b.WriteString("</tbody></table>\n")

	b.WriteString("<h2>Manager changes</h2>\n")
	b.WriteString("<table><thead><tr><th>Manager</th><th>Team</th><th>Total changes</th><th>Latest change</th></tr></thead><tbody>\n")

	// Sort summaries by ManagerName / TeamName
	sortedSummaries := make([]models.ManagerChangeSummary, len(summaries))
	copy(sortedSummaries, summaries)
	sort.Slice(sortedSummaries, func(i, j int) bool {
		n1 := sortedSummaries[i].ManagerName
		if n1 == "" {
			n1 = sortedSummaries[i].ManagerID
		}
		n2 := sortedSummaries[j].ManagerName
		if n2 == "" {
			n2 = sortedSummaries[j].ManagerID
		}
		return n1 < n2
	})

	for _, summary := range sortedSummaries {
		mgrName := summary.ManagerName
		if mgrName == "" {
			mgrName = summary.ManagerID
		}
		teamName := summary.TeamName
		if teamName == "" {
			teamName = mgrName
		}
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td></tr>\n", mgrName, teamName, summary.TotalChanges, summary.LatestChangeAt))
	}
	b.WriteString("</tbody></table>\n")

	b.WriteString("<h2>Historical trends</h2>\n")
	// Sort history by player name / ID
	histKeys := make([]string, 0, len(history))
	for k := range history {
		histKeys = append(histKeys, k)
	}
	sort.Strings(histKeys)
	for _, playerID := range histKeys {
		values := history[playerID]
		if len(values) == 0 {
			continue
		}
		trendStr := ""
		for idx, v := range values {
			if idx > 0 {
				trendStr += ", "
			}
			trendStr += fmt.Sprintf("%s: %d", v.CapturedAt, v.ManagerCount)
		}
		displayName := playerID
		if rec, ok := current[playerID]; ok && rec.PlayerName != "" {
			displayName = rec.PlayerName
		}
		b.WriteString(fmt.Sprintf("<p><strong>%s</strong>: %s</p>\n", displayName, trendStr))
	}
	b.WriteString("</body></html>\n")
	return b.String(), nil
}

// BuildHTMLFromDoc renders static HTML directly from a PlayerOwnershipDoc and change summaries.
func BuildHTMLFromDoc(doc models.PlayerOwnershipDoc, summaries []models.ManagerChangeSummary) (string, error) {
	current := make(map[string]models.PlayerOwnershipRecord, len(doc.Players))
	history := make(map[string][]models.PlayerOwnershipHistoryEntry, len(doc.Players))

	for _, p := range doc.Players {
		current[p.PlayerID] = models.PlayerOwnershipRecord{
			PlayerID:     p.PlayerID,
			PlayerName:   p.PlayerName,
			ManagerCount: p.CurrentCount,
		}
		history[p.PlayerID] = p.History
	}

	return BuildHTML(current, history, summaries)
}
