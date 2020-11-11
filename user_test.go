package destinyhome

import "testing"

func TestGetUser(t *testing.T) {

	// Call the function.
	_, err := getUser("Sydney")
	if err != nil {
		t.Fatalf("couldn't get user: %v", err)
	}
}
