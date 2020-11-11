package destinyhome

import (
	"testing"
)

func TestGetEquipedItem(t *testing.T) {

	// Call the function.
	_, err := getEquipedItem("Sydney", "head", "1")
	if err != nil {
		t.Fatalf("couldn't get item: %v", err)
	}
}
