package destinyhome

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Entry defines a log entry.
type Entry struct {
	Message   string `json:"message"`
	Severity  string `json:"severity,omitempty"`
	Operation string `json:"operation"`
	Trace     string `json:"logging.googleapis.com/trace,omitempty"`

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

func Webhook(w http.ResponseWriter, r *http.Request) {

	var (
		req   webhookRequest
		res   webhookResponse
		trace string
	)

	// Derive the traceID associated with the current request.
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
	}

	// Decode request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Process the request.
	if err := processRequest(&req, &res); err != nil {

		switch errors.Cause(err) {

		case errUserNotFound:
			res.Prompt.FirstSimple.Speech = errors.Cause(err).Error()
			w.WriteHeader(http.StatusNotFound)

		case errCouldntFindItem:
			res.Prompt.FirstSimple.Speech = errors.Cause(err).Error()
			w.WriteHeader(http.StatusNotFound)

		case errOnlyOneAllowed:
			res.Prompt.FirstSimple.Speech = errors.Cause(err).Error()
			w.WriteHeader(http.StatusNotFound)

		case errLoadoutNameTaken:
			res.Prompt.FirstSimple.Speech = errors.Cause(err).Error()
			w.WriteHeader(http.StatusNotFound)

		default:

			log.Println(Entry{
				Message:   errors.Cause(err).Error(),
				Severity:  "CRITICAL",
				Operation: err.Error(),
				Trace:     trace,
			})

			res.Prompt.FirstSimple.Speech = `My backend systems are not working right now, try again later.`
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// Turn the response into bytes.
	js, err := json.Marshal(res)
	if err != nil {

		log.Println(Entry{
			Message:  err.Error(),
			Severity: "CRITICAL",
			Trace:    trace,
		})

		return
	}

	// Set headers.
	r.Header.Set("Content-Type", "application/json")

	// Respond.
	_, err = w.Write(js)
	if err != nil {
	}
}
