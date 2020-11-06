package destinyhome

type modelUser struct {
	username       string
	gamertag       string
	membershipType string
	membershipID   string
	characters     [3]modelCharacter
}

type modelCharacter struct {
	id       string
	loadouts [3]modelLoadout
}

type modelLoadout struct {
	type_ string

	subclassID string

	headID      string
	armsID      string
	chestID     string
	legsID      string
	classItemID string

	kineticID string
	specialID string
	heavyID   string
}
