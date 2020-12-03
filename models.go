package destinyhome

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	AccessToken  string
	RefreshToken string
}

// refresh refreshes the users access token.
func (m *modelUser) refresh() error {

	// Create the url.
	url_ := "https://www.bungie.net/platform/app/oauth/token"

	// Create request body.
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", m.RefreshToken)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	// Create request.
	req, err := http.NewRequest("POST", url_, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Add request headers.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request.
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the status code.
	if res.StatusCode != 200 {
		return errors.New("not 200, got " + http.StatusText(res.StatusCode))
	}

	// Decode the response body.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	x := map[string]interface{}{}

	// Put response into a map.
	if err := json.Unmarshal(body, &x); err != nil {
		return err
	}

	// Set the new access and refresh tokens.
	m.AccessToken = x["access_token"].(string)
	m.RefreshToken = x["refresh_token"].(string)

	return nil
}

func (m *modelUser) save() error {

	// Create context for datastore.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(9)*time.Second,
	)
	defer cancel()

	// Create datastore client.
	c, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	// Update the user.
	key := datastore.IncompleteKey("User", nil)
	if _, err := c.Put(ctx, key, m); err != nil {
		return err
	}

	return nil
}

type modelCharacter struct {
	user *modelUser

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

// getBucket returns the items in a bucket of this character.
func (m *modelCharacter) getBucket(s *bungo.Service, bucket string) ([]string, error) {

	// Get the Character.
	call := s.Destiny2.GetCharacter(m.user.MembershipType, m.user.MembershipID, m.ID).
		Components("201")
	call.Header().Add("Authorization", "Bearer "+m.user.AccessToken)
	res, err := call.Do()
	if err != nil {

		if err != bungo.ErrUnauthorized {
			return nil, err
		}

		// Refresh the access token.
		if err := m.user.refresh(); err != nil {
			return nil, err
		}

		// Save the user now that the token has been refreshed.
		if err := m.user.save(); err != nil {
			return nil, err
		}

		// Try to get the character again.
		call := s.Destiny2.GetCharacter(m.user.MembershipType, m.user.MembershipID, m.ID).
			Components("201")
		call.Header().Add("Authorization", "Bearer "+m.user.AccessToken)
		res, err = call.Do()
		if err != nil {
			return nil, err
		}
	}

	// Create a hash map for the bucket hashes.
	hashMap := map[string]int{
		"kinetic":    1498876634,
		"energy":     2465295065,
		"power":      953998645,
		"head":       3448274439,
		"arms":       3551918588,
		"chest":      14239492,
		"legs":       20886954,
		"class item": 1585787867,
	}

	// Conver the given bucket into its hash.
	bucketHash, ok := hashMap[bucket]
	if !ok {
		return nil, errors.New("bad bucket")
	}

	// Range over the characters items.
	var ret []string
	for _, v := range res.CharacterInventories.Response.Inventory.Data.Items {
		if v.BucketHash != bucketHash {
			continue
		}

		res, err := s.Destiny2.GetDestinyEntityDefinition(
			"DestinyInventoryItemDefinition",
			strconv.Itoa(v.ItemHash),
		).Do()

		if err != nil {
			return nil, err
		}

		ret = append(ret, res.Response.DisplayProperties.Name)
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
	user.Characters[0].user = &user
	user.Characters[1].user = &user
	user.Characters[2].user = &user
	return &user, nil
}
