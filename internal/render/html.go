package render

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

// BuildHTML renders a static HTML summary from current ownership and change summaries.
func BuildHTML(current map[string]models.PlayerOwnershipRecord, history map[string][]models.PlayerOwnershipHistoryEntry, summaries []models.ManagerChangeSummary) (string, error) {
	var b bytes.Buffer
	b.WriteString("<!doctype html>")
	b.WriteString("<html><head><meta charset=\"utf-8\"><title>Guy Sports Team Report</title>")
	b.WriteString("<style>body{font-family:sans-serif;max-width:1100px;margin:2rem auto;padding:0 1rem}table{border-collapse:collapse;width:100%;margin:1rem 0}th,td{border:1px solid #ccc;padding:.5rem;text-align:left}th{background:#f5f5f5}</style></head><body>")
	b.WriteString("<h1>Guy Sports Team Report</h1>")
	b.WriteString("<h2>Player ownership</h2>")
	b.WriteString("<table><thead><tr><th>Player</th><th>Managers</th></tr></thead><tbody>")
	keys := make([]string, 0, len(current))
	for k := range current {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		rec := current[k]
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", rec.PlayerName, rec.ManagerCount))
	}
	b.WriteString("</tbody></table>")
	b.WriteString("<h2>Manager changes</h2>")
	b.WriteString("<table><thead><tr><th>Manager</th><th>Total changes</th><th>Latest change</th></tr></thead><tbody>")
	for _, summary := range summaries {
		b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>", summary.ManagerID, summary.TotalChanges, summary.LatestChangeAt))
	}
	b.WriteString("</tbody></table>")
	b.WriteString("<h2>Historical trends</h2>")
	for playerID, values := range history {
		if len(values) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("<p><strong>%s</strong>: %s</p>", playerID, strings.Trim(strings.Replace(fmt.Sprint(values), "{", "", -1), "}")))
	}
	b.WriteString("</body></html>")
	return b.String(), nil
}
