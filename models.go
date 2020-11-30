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
	Characters     []*modelCharacter
}

type modelCharacter struct {
	user modelUser

	ID       string
	Loadouts []map[string]*modelItem
}

// getCurrentLoadout returns the characters currently equiped loadout.
func (m *modelCharacter) getCurrentLoadout(s *bungo.Service) (map[string]*modelItem, error) {

	// Use bungo to get the character equipment.
	call := s.Destiny2.GetCharacter(m.user.MembershipType, m.user.MembershipID, m.ID)
	res, err := call.Components("CharacterEquipment").Do()
	if err != nil {
		return nil, err
	}

	// Convert the response into a modelLoadout.
	ret := map[string]*modelItem{}
	for _, v := range res.CharacterEquipment.Response.Equipment.Data.Items {

		switch v.BucketHash {

		case 1498876634: // <-- Kinetic.
			ret["kinetic"] = &modelItem{}
			ret["kinetic"].ItemHash = v.ItemHash
			ret["kinetic"].ItemInstanceID = v.ItemInstanceID
		case 2465295065: // <-- Energy.
			ret["energy"] = &modelItem{}
			ret["energy"].ItemHash = v.ItemHash
			ret["energy"].ItemInstanceID = v.ItemInstanceID
		case 953998645: // <-- Power.
			ret["power"] = &modelItem{}
			ret["power"].ItemHash = v.ItemHash
			ret["power"].ItemInstanceID = v.ItemInstanceID
		case 3448274439: // <-- Head.
			ret["head"] = &modelItem{}
			ret["head"].ItemHash = v.ItemHash
			ret["head"].ItemInstanceID = v.ItemInstanceID
		case 3551918588: // <-- Arms.
			ret["arms"] = &modelItem{}
			ret["arms"].ItemHash = v.ItemHash
			ret["arms"].ItemInstanceID = v.ItemInstanceID
		case 14239492: // <-- Chest.
			ret["chest"] = &modelItem{}
			ret["chest"].ItemHash = v.ItemHash
			ret["chest"].ItemInstanceID = v.ItemInstanceID
		case 20886954: // <-- Legs.
			ret["legs"] = &modelItem{}
			ret["legs"].ItemHash = v.ItemHash
			ret["legs"].ItemInstanceID = v.ItemInstanceID
		case 1585787867: // <-- Class-Item.
			ret["class item"] = &modelItem{}
			ret["class item"].ItemHash = v.ItemHash
			ret["class item"].ItemInstanceID = v.ItemInstanceID
		}
	}

	return ret, nil
}

type modelItem struct {
	ItemHash       int
	ItemInstanceID string
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

	// Check if a user  was found.
	if len(users) == 0 {
		return nil, errUserNotFound
	}

	// Return the user.
	user := users[0]
	user.Characters[0].user = user
	user.Characters[1].user = user
	user.Characters[2].user = user
	return &users[0], nil
}
