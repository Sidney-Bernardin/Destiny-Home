package destinyhome

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	port = flag.String("port", "", "The port for the create user webhook to be tested on.")
	see  = flag.Bool("see", false, "See the responses of the Webhook test.")
)

func TestAdmin(t *testing.T) {

	// Handle the CreateUser function.
	http.HandleFunc("/"+adminEndpoint, Admin)

	// Start the server.
	fmt.Printf("Test server for Admin is lintening on :%s...", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		t.Fatalf("Server crashed: %v", err)
	}
}

func TestWebhook(t *testing.T) {

	// Setup test cases.
	tables := []struct {
		handler string
		params  map[string]param
	}{
		{
			handler: "get_equiped_item",
			params: map[string]param{
				"username":       {Resolved: "Sydney"},
				"bucket":         {Resolved: "head"},
				"guardian_index": {Resolved: "1"},
			},
		},
		{
			handler: "get_equiped_item",
			params: map[string]param{
				"username":       {Resolved: "Sydney"},
				"bucket":         {Resolved: "chest"},
				"guardian_index": {Resolved: "0"},
			},
		},
		{
			handler: "get_equiped_item",
			params: map[string]param{
				"username":       {Resolved: "Sydney"},
				"bucket":         {Resolved: "legs"},
				"guardian_index": {Resolved: "2"},
			},
		},
		{
			handler: "get_equiped_item",
			params: map[string]param{
				"username":       {Resolved: "Sydney"},
				"bucket":         {Resolved: "arms"},
				"guardian_index": {Resolved: "1"},
			},
		},
	}

	// Run test cases.
	for _, table := range tables {

		// Print request.
		if *see {
			fmt.Printf("running test case\n\trequest: %v\n", table)
		}

		// Setup request body.
		var req webhookRequest
		req.Handler.Name = table.handler
		req.Intent.Params = table.params

		// Marshal request body.
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("couldn't marshal request: %v", err)
		}

		// Setup http request.
		r, err := http.NewRequest("POST", "localhost:8080/webkooktest", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("couldn't create request: %v", err)
		}

		// Create recorder.
		rec := httptest.NewRecorder()

		// Run the function.
		Webhook(rec, r)

		// Get the response.
		res := rec.Result()

		// Check the status code.
		if res.StatusCode != 200 {
			t.Errorf("status code is not 200, got %s", res.Status)
		}

		// Print response.
		if *see {

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("couldn't read response body: %v", err)
			}

			fmt.Printf("\tresponse: %s\n", string(b))
			fmt.Println("ending test case")
		}
	}
}
