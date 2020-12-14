package destinyhome

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sidney-Bernardin/bungo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pkg/errors"
)

func processRequest(req *webhookRequest, res *webhookResponse) error {

	const operation = "processRequest"

	switch req.Handler.Name {

	case "get_equiped_item":

		// Get params.
		username := req.Intent.Params["username"].Resolved
		bucket := req.Intent.Params["bucket"].Resolved
		guardianIndex := req.Intent.Params["guardian_index"].Resolved

		// Get the item.
		item, err := handleGetEquipedItem(username, bucket, guardianIndex)
		if err != nil {
			return errors.Wrap(err, operation+": handle get equip item failed")
		}

		// Setup the response.
		res.Prompt.FirstSimple.Speech = fmt.Sprintf("Your %s is %s", bucket, item)

	case "equip_item":

		// Get params.
		username := req.Intent.Params["username"].Resolved
		guardianIndex := req.Intent.Params["guardian_index"].Resolved
		itemName := req.Intent.Params["item_name"].Resolved

		// Equip the item.
		if err := handleEquipItem(username, guardianIndex, itemName); err != nil {
			return errors.Wrap(err, operation+": handle equip item failed")
		}

		// Setup the response.
		res.Prompt.FirstSimple.Speech = fmt.Sprintf("Done equiping %s!", itemName)

	case "save_loadout":

		// Get params.
		username := req.Intent.Params["username"].Resolved
		guardianIndex := req.Intent.Params["guardian_index"].Resolved
		loadoutName := req.Intent.Params["loadout_name"].Resolved

		// Equip the item.
		if err := saveLoadout(username, guardianIndex, loadoutName); err != nil {
			return errors.Wrap(err, operation+": handle save loadout item failed")
		}

		// Setup the response.
		res.Prompt.FirstSimple.Speech = fmt.Sprintf("Saved %s!", loadoutName)
	}

	// Setup common response fields.
	res.Session.ID = req.Session.ID
	res.Scene.Name = req.Scene.Name

	return nil
}

func handleGetEquipedItem(username, bucket, guardianIndex string) (string, error) {
	const operation = "handleGetEquipedItem"

	// Get the user given the username.
	user, err := getUser(username)
	if err != nil {
		return "", errors.Wrap(err, operation+": getting user failed")
	}

	// Create bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		return "", errors.Wrap(err, operation+": new bungo service failed")
	}

	// Convert the guardianIndex into an int.
	number, err := strconv.Atoi(guardianIndex)
	if err != nil {
		return "", errors.Wrap(err, operation+": string to int failed")
	}

	// Get the currently equiped loadout.
	res, err := user.Characters[number].getCurrentLoadout(s)
	if err != nil {
		return "", errors.Wrap(err, operation+": getting loadout failed")
	}

	// Get the item definition of the item hash.
	res2, err := s.Destiny2.GetDestinyEntityDefinition(
		"DestinyInventoryItemDefinition", strconv.Itoa(res[bucket].ItemHash)).Do()
	if err != nil {
		return "", errors.Wrap(err, operation+": getting entity definition failed")
	}

	// Return the name of the item.
	return res2.Response.DisplayProperties.Name, nil
}

func handleEquipItem(username, guardianIndex, itemName string) error {

	const operation = "handleEquipItem"

	// Get the user given the username.
	user, err := getUser(username)
	if err != nil {
		return errors.Wrap(err, operation+": get user failed")
	}

	// Create bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		return errors.Wrap(err, operation+": new bungo service failed")
	}

	// Convert the guardianIndex into an int.
	number, err := strconv.Atoi(guardianIndex)
	if err != nil {
		return errors.Wrap(err, operation+": string to int failed")
	}

	successChan := make(chan struct{})
	errChan := make(chan error)
	notFoundChan := make(chan struct{})

	for k := range gearHashMap {

		go func(k string) {

			// Get the items from the bucket.
			bucket, err := user.Characters[number].getBucket(s, k)
			if err != nil {
				errChan <- errors.Wrap(err, operation+": getting bucket failed")
				return
			}

			// Convert bucket into a slice of the item names.
			names := []string{}
			for _, v := range bucket {
				names = append(names, strings.ToLower(v.Name))
			}

			// Fuzzy search for the item.
			res := fuzzy.Find(strings.ToLower(itemName), names)
			if len(res) == 0 {
				notFoundChan <- struct{}{}
				return
			}

			// Equip the item.
			for i := range bucket {
				if strings.ToLower(bucket[i].Name) == res[0] {

					err = user.Characters[number].equipItem(s, bucket[i].ItemInstanceID)
					if err != nil {
						errChan <- errors.Wrap(err, operation+": equip item failed")
						return
					}

					if err := user.save(); err != nil {
						errChan <- errors.Wrap(err, operation+": save user failed")
						return
					}

					successChan <- struct{}{}
				}
			}
		}(k)
	}

	for range gearHashMap {
		select {
		case _ = <-successChan:
			return nil
		case err := <-errChan:
			return err
		case _ = <-notFoundChan:
			continue
		}
	}

	return errCouldntFindItem
}

func saveLoadout(username, guardianIndex, loadoutName string) error {

	const operation = "saveLoadout"

	// Get the user given the username.
	user, err := getUser(username)
	if err != nil {
		return errors.Wrap(err, operation+": get user failed")
	}

	// Create bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		return errors.Wrap(err, operation+": new bungo service failed")
	}

	// Convert the guardianIndex into an int.
	number, err := strconv.Atoi(guardianIndex)
	if err != nil {
		return errors.Wrap(err, operation+": string to int failed")
	}

	// Check if the loadoutName is unique.
	for _, v := range user.Characters[number].Loadouts {
		if v.Name == loadoutName {
			return errLoadoutNameTaken
		}
	}

	// Get current loadout.
	currentLoadout, err := user.Characters[number].getCurrentLoadout(s)
	if err != nil {
		return errors.Wrap(err, operation+":get current loadout failed")
	}

	// Add the loadout to the user.
	user.Characters[number].Loadouts = append(user.Characters[number].Loadouts, modelLoadout{
		Name:      loadoutName,
		Kinetic:   currentLoadout["kinetic"],
		Special:   currentLoadout["special"],
		Power:     currentLoadout["power"],
		Head:      currentLoadout["head"],
		Arms:      currentLoadout["arms"],
		Chest:     currentLoadout["chest"],
		Legs:      currentLoadout["legs"],
		ClassItem: currentLoadout["class item"],
	})

	// Save the user.
	if err := user.save(); err != nil {
		return errors.Wrap(err, operation+":save user failed")
	}

	return nil
}
