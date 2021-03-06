package destinyhome

import "flag"

var (
	apiKey       string
	projectID    string
	clientID     string
	clientSecret string

	see                = flag.Bool("see", false, "See the responses of the unit tests.")
	testUsername       string
	testMemType        string
	testMemID          string
	testItemName       string
	testItemID         string
	testAccessToken    string
	testRefreshToken   string
	testCharacterID    string
	testCharacterIndex string
	testLoadoutName    string

	gearHashMap = map[string]int{
		"kinetic":    1498876634,
		"special":    2465295065,
		"power":      953998645,
		"head":       3448274439,
		"arms":       3551918588,
		"chest":      14239492,
		"legs":       20886954,
		"class_item": 1585787867,
	}
)
