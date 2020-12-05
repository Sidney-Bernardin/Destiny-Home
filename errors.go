package destinyhome

import "errors"

var (
	errUserNotFound    = errors.New("user not found")
	errCouldntFindItem = errors.New("could not find item")
)
