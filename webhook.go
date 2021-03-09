package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/logging"
	"github.com/pkg/errors"
)

// Entry defines a log entry.
type Entry struct {
	Message   string `json:"message"`
	Severity  string `json:"severity,omitempty"`
	Operation string `json:"operation"`

	// Cloud Log Viewer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) String() string {

	if e.Severity == "" {
		e.Severity = "INFO"
	}

	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}

	return string(out)
}

type webhookRequest struct {
	Intent struct {
		Name   string `json:"name"`
		Params map[string]struct {
			Original string `json:"original"`
			Resolved string `json:"resolved"`
		} `json:"params"`
		Query string `json:"query"`
	} `json:"intent"`
}

type webhookResponse struct {
	Prompt struct {
		FirstSimple struct {
			Speech string `json:"speech"`
			Text   string `json:"text"`
		} `json:"firstSimple"`
		Override bool `json:"override"`
	} `json:"prompt"`
	Scene struct {
		Name string `json:"name"`
		Next struct {
			Name string `json:"name"`
		} `json:"next"`
		Slots struct{} `json:"slots"`
	} `json:"scene"`
	Session struct {
		ID     string   `json:"id"`
		Params struct{} `json:"params"`
	} `json:"session"`
}

func init() {

}

func Webhook(w http.ResponseWriter, r *http.Request) {

	var (
		req *webhookRequest
		res *webhookResponse
	)

	// Create logger.
	ctx := context.Background()
	client, err := logging.NewClient(ctx, "my-project")
	if err != nil {
		log.Println(logging.Entry{
			Severity: logging.Critical,
			Payload:  errors.Wrap(err, "could not create logging client").Error(),
		})
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Decode request body.
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Process the request.
	res, err := doRequest(req)
	if err != nil {

	}
}

func doRequest(req *webhookRequest) (*webhookResponse, error) {
	return nil, nil
}
