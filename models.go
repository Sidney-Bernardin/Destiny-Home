package destinyhome

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Sidney-Bernardin/bungo"
)

type modelUser struct {
	Username       string
	Gamertag       string
	MembershipType string
	MembershipID   string
	Characters     []modelCharacter
}

type modelCharacter struct {
	user modelUser

	ID       string
	Loadouts []modelLoadout
}

// TODO: Finish this function.
// getCurrentLoadout returns the characters currently equiped loadout.
func (m *modelCharacter) getCurrentLoadout(s *bungo.Service) (*modelLoadout, error) {

	call := s.Destiny2.GetCharacter(m.user.MembershipType, m.user.MembershipID, m.ID)
	res, err := call.Components("CharacterEquipment").Do()
	if err != nil {
		return nil, err
	}

	bucketHashTable := map[int]string{
		1498876634: "kinetic",
		2465295065: "energy",
		953998645:  "power",
		3448274439: "head",
		3551918588: "arms",
		14239492:   "chest",
		1585787867: "class-item",
	}

	var ret *modelLoadout

	for _, v := range res.CharacterEquipment.Response.Equipment.Data.Items {
		if _, ok := bucketHashTable[v.BucketHash]; !ok {
			continue
		}
	}

	return ret, nil
}

type modelLoadout struct {
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

func (m *modelCharacter) save() {

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

	if len(users) == 0 {
		return nil, errUserNotFound
	}

	return &users[0], nil
}
