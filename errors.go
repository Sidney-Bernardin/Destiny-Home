package destinyhome

import "errors"

var (
	errUserNotFound     = errors.New("user not found")
	errCouldntFindItem  = errors.New("could not find item")
	errOnlyOneAllowed   = errors.New("you can only have one item of this type equipped")
	errLoadoutNameTaken = errors.New("loadout name already exists")
)
