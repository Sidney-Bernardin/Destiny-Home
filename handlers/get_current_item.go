package main

import (
	"net/http"
	"strconv"

	"github.com/Sidney-Bernardin/bungo"
	"github.com/pkg/errors"
)

func GetCurrentItem(username, guardianIndex, itemName string) error {

	const operation = ""

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
