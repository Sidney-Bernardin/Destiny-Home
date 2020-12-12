package destinyhome

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Sidney-Bernardin/bungo"
	"github.com/pkg/errors"
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

	const operation = "modelUser.refresh"

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
		return errors.Wrap(err, operation+": creating request failed")
	}

	// Add request headers.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request.
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, operation+": doing request failed")
	}
	defer res.Body.Close()

	// Check the status code.
	if res.StatusCode != 200 {
		err := errors.New(http.StatusText(res.StatusCode))
		return errors.Wrap(err, operation+": not 200")
	}

	// Decode the response body.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	x := map[string]interface{}{}

	// Put response into a map.
	if err := json.Unmarshal(body, &x); err != nil {
		return errors.Wrap(err, operation+": unmarshal failed")
	}

	// Set the new access and refresh tokens.
	m.AccessToken = x["access_token"].(string)
	m.RefreshToken = x["refresh_token"].(string)

	return nil
}

// save updates the user in the database.
func (m *modelUser) save() error {

	const operation = "modelUser.save"

	// Create context for datastore.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(9)*time.Second,
	)
	defer cancel()

	// Create datastore client.
	c, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, operation+": new datastore client failed")
	}

	// Update the user.
	key := datastore.NameKey("User", m.Username, nil)
	if _, err := c.Put(ctx, key, m); err != nil {
		return errors.Wrap(err, operation+": put user failed")
	}

	return nil
}

type modelCharacter struct {
	user *modelUser

	ID       string
	Loadouts []map[string]*modelItem
}

// getCurrentLoadout returns the characters currently equiped loadout.
func (m *modelCharacter) getCurrentLoadout(s *bungo.Service) (map[string]modelItem, error) {

	const operation = "modelCharacter.getCurrentLoadout"

	// Use bungo to get the character equipment.
	call := s.Destiny2.GetCharacter(m.user.MembershipType, m.user.MembershipID, m.ID)
	res, err := call.Components("CharacterEquipment").Do()
	if err != nil {
		return nil, errors.Wrap(err, operation+": get character failed")
	}

	// range over the response items.
	ret := map[string]modelItem{}
	for _, v := range res.CharacterEquipment.Response.Equipment.Data.Items {

		// Loop over the grear hash map.
		for bucketName, bucketHash := range gearHashMap {
			if v.BucketHash == bucketHash {

				// Set the item.
				ret[bucketName] = modelItem{
					ItemHash:       v.ItemHash,
					ItemInstanceID: v.ItemInstanceID,
				}
			}
		}
	}

	return ret, nil
}

// getBucket returns the items in a bucket of this character.
func (m *modelCharacter) getBucket(s *bungo.Service, bucket string) ([]modelItem, error) {

	const operation = "modelCharacter.getBucket"

	// Get the Character.
	call := s.Destiny2.GetCharacter(m.user.MembershipType, m.user.MembershipID, m.ID).
		Components("201")
	call.Header().Add("Authorization", "Bearer "+m.user.AccessToken)
	res, err := call.Do()
	if err != nil {

		if err != bungo.ErrUnauthorized {
			return nil, errors.Wrap(err, operation+": get character failed")
		}

		// Refresh the access token.
		if err := m.user.refresh(); err != nil {
			return nil, errors.Wrap(err, operation+": refreshing token failed")
		}

		// Try to get the character again.
		call := s.Destiny2.GetCharacter(m.user.MembershipType, m.user.MembershipID, m.ID).
			Components("201")
		call.Header().Add("Authorization", "Bearer "+m.user.AccessToken)
		res, err = call.Do()
		if err != nil {
			return nil, errors.Wrap(err, operation+": get character failed")
		}
	}

	// Convert the given bucket into its hash.
	bucketHash, ok := gearHashMap[bucket]
	if !ok {
		err := errors.New("bad bucket")
		return nil, errors.Wrap(err, operation+": convert to hash failed")
	}

	errChan := make(chan error)
	continueChan := make(chan struct{})

	// Range over the characters items.
	var ret []modelItem
	for _, v := range res.CharacterInventories.Response.Inventory.Data.Items {

		go func(v bungo.SingleComponentResponseOfDestinyInventoryComponentItem) {

			if v.BucketHash != bucketHash {
				continueChan <- struct{}{}
				return
			}

			// Get the item's entity definition to get its name.
			res, err := s.Destiny2.GetDestinyEntityDefinition(
				"DestinyInventoryItemDefinition",
				strconv.Itoa(v.ItemHash),
			).Do()

			if err != nil {
				errChan <- errors.Wrap(err, operation+": get item definition failed")
				return
			}

			// Add the item to ret.
			ret = append(ret, modelItem{
				Name:           res.Response.DisplayProperties.Name,
				ItemHash:       v.ItemHash,
				ItemInstanceID: v.ItemInstanceID,
			})

			continueChan <- struct{}{}
		}(v)
	}

	for range res.CharacterInventories.Response.Inventory.Data.Items {
		select {
		case err := <-errChan:
			return nil, err
		case _ = <-continueChan:
			continue
		}
	}

	return ret, nil
}

func (m *modelCharacter) equipItem(s *bungo.Service, itemInstanceID string) error {

	const operation = "modelCharacter.equipItem"

	// Create the request body.
	data := &bungo.ItemActionRequest{
		ItemId:         itemInstanceID,
		CharacterId:    m.ID,
		MembershipType: m.user.MembershipType,
	}

	// Equip the item.
	call := s.Destiny2.EquipItem(data)
	call.Header().Add("Authorization", "Bearer "+m.user.AccessToken)
	_, err := call.Do()
	if err != nil {

		if err != bungo.ErrUnauthorized {

			if _, ok := err.(*bungo.BungoError); ok {
				if err.(*bungo.BungoError).ErrorCode == 1641 {
					return errOnlyOneAllowed
				}
			}

			return errors.Wrap(err, operation+": equip item failed")
		}

		// Refresh the access token.
		if err := m.user.refresh(); err != nil {
			return errors.Wrap(err, operation+": user refresh failed")
		}

		// Try to equip the item again.
		call := s.Destiny2.EquipItem(data)
		call.Header().Add("Authorization", "Bearer "+m.user.AccessToken)
		_, err = call.Do()
		if err != nil {

			if _, ok := err.(*bungo.BungoError); ok {
				if err.(*bungo.BungoError).ErrorCode == 1641 {
					return errOnlyOneAllowed
				}
			}

			return errors.Wrap(err, operation+": equip item failed after refresh")
		}
	}

	return nil
}

type modelItem struct {
	Name           string
	ItemHash       int
	ItemInstanceID string
}

// getUser returns a user from the database given a username.
func getUser(username string) (*modelUser, error) {

	const operation = "getUser"

	// Create context for datastore.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(9)*time.Second,
	)
	defer cancel()

	// Create datastore client.
	c, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, operation+": new datastore client failed")
	}

	// Create a query to get the user..
	query := datastore.NewQuery("User").
		Filter("Username =", username)

	// Run the query.
	var users []modelUser
	_, err = c.GetAll(ctx, query, &users)
	if err != nil {
		return nil, errors.Wrap(err, operation+": get all users failed")
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
