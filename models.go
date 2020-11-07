package destinyhome

type modelUser struct {
	Username       string
	Gamertag       string
	MembershipType string
	MembershipID   string
	Characters     []modelCharacter
}

type modelCharacter struct {
	ID       string
	Loadouts []modelLoadout
}

type modelLoadout struct {
	Type_ string

	SubclassID string

	HeadID      string
	ArmsID      string
	ChestID     string
	LegsID      string
	ClassItemID string

	KineticID string
	SpecialID string
	HeavyID   string
}
