package unit

import (
	"testing"

	"github.com/thisisdavidbell/ff-sdd/internal/capture"
)

func TestExtractManagersFromHTMLAcceptsTeamDetailLinkVariants(t *testing.T) {
	page := `<a href="javascript:MM_openBrWindow('publicteamdetail.php?id=8&amp;name=Finsbury+Bridge&amp;view=public')">Team</a>
<a href="publicteamdetail.php?name=Kopper+Highlights&amp;view=public&amp;id=20">Team</a>
<a onclick="MM_openBrWindow('publicteamdetail.php?id=31&name=Chopper%20Squad&view=public')">Team</a>`

	managers := capture.ExtractManagersFromHTML(page)
	if len(managers) != 3 {
		t.Fatalf("expected 3 managers, got %d: %#v", len(managers), managers)
	}

	if managers[0].ID != "8" || managers[0].Name != "Finsbury Bridge" {
		t.Fatalf("unexpected first manager: %#v", managers[0])
	}
	if managers[1].ID != "20" || managers[1].Name != "Kopper Highlights" {
		t.Fatalf("unexpected second manager: %#v", managers[1])
	}
	if managers[2].ID != "31" || managers[2].Name != "Chopper Squad" {
		t.Fatalf("unexpected third manager: %#v", managers[2])
	}
}
