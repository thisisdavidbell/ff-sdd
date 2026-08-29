package unit

import (
	"testing"

	"github.com/thisisdavidbell/ff-sdd/internal/capture"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
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

func TestFormatWithUnderscores(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Finsbury Bridge", "Finsbury_Bridge"},
		{"  Willy's   Wanderers  ", "Willy's_Wanderers"},
		{"SingleName", "SingleName"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := capture.FormatWithUnderscores(tc.in); got != tc.want {
			t.Errorf("FormatWithUnderscores(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsRosterUnchanged(t *testing.T) {
	r1 := []models.PlayerReference{
		{PlayerID: "p1"},
		{PlayerID: "p2"},
	}
	r2 := []models.PlayerReference{
		{PlayerID: "p2"},
		{PlayerID: "p1"},
	}
	r3 := []models.PlayerReference{
		{PlayerID: "p1"},
		{PlayerID: "p3"},
	}

	if !capture.IsRosterUnchanged(r1, r2) {
		t.Fatal("expected r1 and r2 to be unchanged")
	}
	if capture.IsRosterUnchanged(r1, r3) {
		t.Fatal("expected r1 and r3 to be changed")
	}
}
