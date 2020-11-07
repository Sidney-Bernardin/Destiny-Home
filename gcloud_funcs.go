package destinyhome

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Sidney-Bernardin/bungo"
)

// CreateUser is a Google Cloud Function for creating users.
// After it asks the user for their destiny info, it adds them
// to a Google Cloud Datastore, Database.
func CreateUser(w http.ResponseWriter, r *http.Request) {

	// If form was submited (if request method is POST), create and store the
	// user into the database.
	if r.Method == "POST" {

		// Get form values.
		username := r.FormValue("username")
		membershipType := r.FormValue("membershipType")
		gamertag := r.FormValue("gamertag")

		// Create bungo service.
		s, err := bungo.NewService(&http.Client{}, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Use SearchDestinyPlayer to get the membershipID.
		res, err := s.Destiny2.SearchDestinyPlayer(membershipType, gamertag).Do()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the destiny player was found.
		if len(res.Response) == 0 {

			err := "couldn't find destiny player: " + gamertag
			http.Error(w, err, http.StatusInternalServerError)
			return
		}

		// Extract the membershipID from the response.
		membershipID := res.Response[0].MembershipID

		// Use Profile to get the characterIDs.
		res2, err := s.Destiny2.GetProfile(membershipType, membershipID).
			Components("Profiles").Do()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Extract the characterIDs from the response.
		characterIDs := res2.Profiles.Response.Profile.Data.CharacterIds
		switch len(characterIDs) {
		case 0:
			characterIDs = []string{"", "", ""}
		case 1:
			characterIDs = []string{characterIDs[0], "", ""}
		case 2:
			characterIDs = []string{characterIDs[0], characterIDs[1], ""}
		case 3:
			characterIDs = []string{characterIDs[0], characterIDs[1], characterIDs[2]}
		}

		// Assemble the user.
		user := modelUser{
			Username:       username,
			Gamertag:       gamertag,
			MembershipType: membershipType,
			MembershipID:   membershipID,
			Characters: []modelCharacter{
				{ID: characterIDs[0]},
				{ID: characterIDs[1]},
				{ID: characterIDs[2]},
			},
		}

		// Create context for datastore.
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()

		// Create a datastore client.
		c, err := datastore.NewClient(ctx, projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store the user into the database.
		key := datastore.IncompleteKey("User", nil)
		if _, err := c.Put(ctx, key, &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	// Execute the, create_user.html template.
	err := temps.ExecuteTemplate(w, "create_user.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MainWebHook(w http.ResponseWriter, r *http.Request) {}
