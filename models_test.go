package destinyhome

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Sidney-Bernardin/bungo"
)

func TestGetCurrentLoadout(t *testing.T) {

	// Require membership type.
	if testMemType == "" || testMemType == " " {
		t.Fatal("TEST_MEM_TYPE must be set")
	}

	// Require membership id.
	if testMemID == "" || testMemID == " " {
		t.Fatal("TEST_MEM_ID must be set")
	}

	// Create the user.
	user := modelUser{
		MembershipType: testMemType,
		MembershipID:   testMemID,
	}

	// Create a character.
	user.Characters = append(user.Characters, &modelCharacter{
		user: &user,

		ID: "2305843009294908257",
	})

	// Create a bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		t.Fatalf("couldn't create bungo service: %v", err)
	}

	// Get current loadout.
	res, err := user.Characters[0].getCurrentLoadout(s)
	if err != nil {
		t.Fatalf("couldn't get loadout: %v", err)
	}

	// Print response.
	if *see {
		fmt.Println(res)
	}
}

func TestGetBucket(t *testing.T) {

	// Require membership type.
	if testMemType == "" || testMemType == " " {
		t.Fatal("TEST_MEM_TYPE must be set")
	}

	// Require membership id.
	if testMemID == "" || testMemID == " " {
		t.Fatal("TEST_MEM_ID must be set")
	}

	// Require access token.
	if testAccessToken == "" || testAccessToken == " " {
		t.Fatal("TEST_ACCESS_TOKEN must be set")
	}

	// Require refresh token.
	if testRefreshToken == "" || testRefreshToken == " " {
		t.Fatal("TEST_REFRESH_TOKEN must be set")
	}

	// Create user.
	user := modelUser{
		MembershipType: testMemType,
		MembershipID:   testMemID,
		AccessToken:    testAccessToken,
		RefreshToken:   testRefreshToken,
	}

	// Create character.
	user.Characters = append(user.Characters, &modelCharacter{
		user: &user,
		ID:   testCharacterID,
	})

	// Create bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		t.Fatalf("couldn't create bungo service: %v", err)
	}

	// Loop over the gear hash map to get all the buckets.
	for bucket := range gearHashMap {

		// Get the bucket.
		res, err := user.Characters[0].getBucket(s, bucket)
		if err != nil {
			t.Fatalf("couldn't get bucket: %v", err)
		}

		// Print the response.
		if *see {
			fmt.Printf("%s: %v\n\n", bucket, res)
		}
	}
}

func TestEquipItem(t *testing.T) {

	// Require membership type.
	if testMemType == "" || testMemType == " " {
		t.Fatal("TEST_MEM_TYPE must be set")
	}

	// Require membership id.
	if testMemID == "" || testMemID == " " {
		t.Fatal("TEST_MEM_ID must be set")
	}

	// Require item id.
	if testItemID == "" || testItemID == " " {
		t.Fatal("TEST_ITEM_ID must be set")
	}

	// Require access token.
	if testAccessToken == "" || testAccessToken == " " {
		t.Fatal("TEST_ACCESS_TOKEN must be set")
	}

	// Require refresh token.
	if testRefreshToken == "" || testRefreshToken == " " {
		t.Fatal("TEST_REFRESH_TOKEN must be set")
	}

	// Create user.
	user := modelUser{
		MembershipType: testMemType,
		MembershipID:   testMemID,
		AccessToken:    testAccessToken,
		RefreshToken:   testRefreshToken,
	}

	// Create character.
	user.Characters = append(user.Characters, &modelCharacter{
		user: &user,
		ID:   testCharacterID,
	})

	// Create bungo service.
	s, err := bungo.NewService(&http.Client{}, apiKey)
	if err != nil {
		t.Fatalf("couldn't create bungo service: %v", err)
	}

	// Get the bucket.
	err = user.Characters[0].equipItem(s, testItemID)
	if err != nil {
		fmt.Println(user.AccessToken)
		t.Fatalf("couldn't equip item: %v", err)
	}
}
