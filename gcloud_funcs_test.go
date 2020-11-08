package destinyhome

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	port = flag.String("port", "", "The port for the create user webhook to be tested on.")
	see  = flag.Bool("see", false, "See the responses of the Webhook test.")
)

func TestCreateUser(t *testing.T) {

	// Handle the CreateUser function.
	http.HandleFunc("/"+createUserEndpoint, CreateUser)

	// Start the server.
	fmt.Printf("Test server for CreateUser is lintening on :%s...", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		t.Fatalf("Server crashed: %v", err)
	}
}

func TestWebhook(t *testing.T) {

	// Setup tast cases.
	tables := []struct {
		handler string
		params  map[string]param
	}{
		{
			handler: "get_item",
			params: map[string]param{
				"account": {
					Resolved: "sidney",
				},
				"item": {
					Resolved: "helmet",
				},
				"guardian_index": {
					Resolved: "1",
				},
			},
		},
	}

	// Run test cases.
	for _, table := range tables {

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
			t.Fatalf("Status code is not 200, got %s", res.Status)
		}

		if *see {

			var res2 webhookResponse
			if err := json.NewDecoder(res.Body).Decode(&res2); err != nil {
				t.Fatalf("couldn't decode reponse: %v", err)
			}

			fmt.Println(res2.Prompt.FirstSimple.Speech)
		}
	}
}
