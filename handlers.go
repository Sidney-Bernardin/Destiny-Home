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
