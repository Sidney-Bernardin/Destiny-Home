package destinyhome

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type modelUser struct {
	Username       string
	Gamertag       string
	MembershipType string
	MembershipID   string
	Characters     []modelCharacter
}

type modelCharacter struct {
	ID       string
	Loadouts []modelLoadout
}

type modelLoadout struct {
	Type_ string

	SubclassID string

	HeadID      string
	ArmsID      string
	ChestID     string
	LegsID      string
	ClassItemID string

	KineticID string
	SpecialID string
	HeavyID   string
}

// getUser returns a user from the database given a username.
func getUser(username string) (*modelUser, error) {

	// Create context for datastore.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(9)*time.Second,
	)
	defer cancel()

	// Create datastore client.
	c, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Create a query to get the user..
	query := datastore.NewQuery("User").
		Filter("Username =", username)

	// Run the query.
	var users []modelUser
	_, err = c.GetAll(ctx, query, &users)
	if err != nil {
		return nil, err
	}

	// TODO: Get rid of this
	if len(users) == 0 {
		return nil, errUserNotFound
	}

	return &users[0], nil
}
