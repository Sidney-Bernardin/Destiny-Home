package destinyhome

import (
	"flag"
	"fmt"
	"net/http"
	"testing"
)

var (
	port = flag.String("port", "", "The port for the create user webhook to be tested on.")
)

func TestCreateUser(t *testing.T) {

	// Handle the CreateUser function.
	http.HandleFunc("/create_user", CreateUser)

	// Start the server.
	fmt.Printf("Test server for CreateUser is lintening on :%s...", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		t.Fatalf("Server crashed: %v", err)
	}
}
