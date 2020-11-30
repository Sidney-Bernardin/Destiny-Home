package destinyhome

import (
	"net/http"
	"testing"

	"github.com/Sidney-Bernardin/bungo"
)

func TestGetCurrentLoadout(t *testing.T) {

	// Create user.
	user := modelUser{
		MembershipType: "2",
		MembershipID:   "4611686018458149036",
	}

	// Create a character for the user.
	user.Characters = append(user.Characters, &modelCharacter{
		user: user,

		ID: "2305843009294908257",
	})

	// Create a new bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		t.Fatalf("couldn't create bungo service: %v", err)
	}

	// Get the Character's current loadout.
	_, err = user.Characters[0].getCurrentLoadout(s)
	if err != nil {
		t.Fatalf("couldn't get loadout: %v", err)
	}
}
