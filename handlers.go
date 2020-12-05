package destinyhome

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sidney-Bernardin/bungo"
)

func processRequest(req *webhookRequest, res *webhookResponse) error {

	switch req.Handler.Name {

	case "get_equiped_item":

		// Get params.
		username := req.Intent.Params["username"].Resolved
		bucket := req.Intent.Params["bucket"].Resolved
		guardianIndex := req.Intent.Params["guardian_index"].Resolved

		// Get the item.
		item, err := getEquipedItem(username, bucket, guardianIndex)
		if err != nil {
			return err
		}

		// Setup the response.
		res.Prompt.FirstSimple.Speech = fmt.Sprintf("Your %s is %s", bucket, item)

	case "equip_item":

		// Get params.
		username := req.Intent.Params["username"].Resolved
		guardianIndex := req.Intent.Params["guardian_index"].Resolved
		itemName := req.Intent.Params["item_name"].Resolved

		// Equip the item.
		if err := equipItem(username, guardianIndex, itemName); err != nil {
			return err
		}

		// Setup the response.
		res.Prompt.FirstSimple.Speech = fmt.Sprintf("Done equiping %s!", itemName)
	}

	// Setup common response fields.
	res.Session.ID = req.Session.ID
	res.Scene.Name = req.Scene.Name

	return nil
}

func getEquipedItem(username, bucket, guardianIndex string) (string, error) {

	// Get the user given the username.
	user, err := getUser(username)
	if err != nil {
		return "", err
	}

	// Create bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		return "", err
	}

	// Convert the guardianIndex into an int.
	number, err := strconv.Atoi(guardianIndex)
	if err != nil {
		return "", err
	}

	// Get the currently equiped loadout.
	res, err := user.Characters[number].getCurrentLoadout(s)
	if err != nil {
		return "", err
	}

	// Get the item definition of the item hash.
	res2, err := s.Destiny2.GetDestinyEntityDefinition(
		"DestinyInventoryItemDefinition", strconv.Itoa(res[bucket].ItemHash)).Do()
	if err != nil {
		return "", err
	}

	// Return the name of the item.
	return res2.Response.DisplayProperties.Name, nil
}

func equipItem(username, guardianIndex, itemName string) error {

	// Get the user given the username.
	user, err := getUser(username)
	if err != nil {
		return err
	}

	// Create bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		return err
	}

	// Convert the guardianIndex into an int.
	number, err := strconv.Atoi(guardianIndex)
	if err != nil {
		return err
	}

	for k := range gearHashMap {

		bucket, err := user.Characters[number].getBucket(s, k)
		if err != nil {
			return err
		}

		err = user.Characters[number].equipItem(s, bucket[0].ItemInstanceID)
		if err != nil {
			return err
		}

		return nil
	}

	return errCouldntFindItem
}
