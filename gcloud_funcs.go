package destinyhome

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

// Entry defines a log entry.
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

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

// Admin responds with a template allowing an administrator to login
// and manage users.
func Admin(w http.ResponseWriter, r *http.Request) {

	var (
		trace string
		users = []modelUser{}
	)

	// Derive the traceID associated with the current request.
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
	}

	switch r.URL.Query()["process"][0] {

	case "link_with_bungie":

		// Authenticate the login.
		password := r.FormValue("password")
		if password != adminPassword {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

	case "home":

		// Authenticate the login.
		password := r.FormValue("password")
		if password != adminPassword {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Create context for datastore.
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()

		// Create datastore client.
		c, err := datastore.NewClient(ctx, projectID)
		if err != nil {

			log.Println(Entry{
				Message:  err.Error(),
				Severity: "CRITICAL",
				Trace:    trace,
			})

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get all user from the database.
		q := datastore.NewQuery("User")
		if _, err := c.GetAll(ctx, q, &users); err != nil {

			log.Println(Entry{
				Message:  err.Error(),
				Severity: "CRITICAL",
				Trace:    trace,
			})

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Execute template and pass the users into it.
	err := temps.ExecuteTemplate(w, "admin.html", map[string]interface{}{
		"Users": users,
	})

	if err != nil {

		log.Println(Entry{
			Message:  err.Error(),
			Severity: "CRITICAL",
			Trace:    trace,
		})

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

		switch err {

		case errUserNotFound:
			res.Prompt.FirstSimple.Speech = `I couldn't find you in my database.`
			w.WriteHeader(http.StatusNotFound)

		default:

			log.Println(Entry{
				Message:  err.Error(),
				Severity: "CRITICAL",
				Trace:    trace,
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
