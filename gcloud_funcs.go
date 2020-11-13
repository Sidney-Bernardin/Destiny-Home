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
	"github.com/Sidney-Bernardin/bungo"
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

// CreateUser is a Google Cloud Function for creating users.
// After it asks the user for their destiny info, it adds them
// to a Google Cloud Datastore, Database.
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var trace string

	// Derive the traceID associated with the current request.
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
	}

	// If form was submited (if request method is POST), create and store the
	// user into the database.
	if r.Method == "POST" {

		// Get form values.
		username := r.FormValue("username")
		membershipType := r.FormValue("membershipType")
		gamertag := r.FormValue("gamertag")

		// Create bungo service.
		s, err := bungo.NewService(&http.Client{}, apiKey)
		if err != nil {

			log.Println(Entry{
				Message:  err.Error(),
				Severity: "CRITICAL",
				Trace:    trace,
			})

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Use SearchDestinyPlayer to get the membershipID.
		res, err := s.Destiny2.SearchDestinyPlayer(membershipType, gamertag).Do()
		if err != nil {

			log.Println(Entry{
				Message:  err.Error(),
				Severity: "CRITICAL",
				Trace:    trace,
			})

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the destiny player was found.
		if len(res.Response) == 0 {

			err := "couldn't find destiny player: " + gamertag
			http.Error(w, err, http.StatusNotFound)
			return
		}

		// Extract the membershipID from the response.
		membershipID := res.Response[0].MembershipID

		// Use Profile to get the characterIDs.
		res2, err := s.Destiny2.GetProfile(membershipType, membershipID).
			Components("Profiles").Do()
		if err != nil {

			log.Println(Entry{
				Message:  err.Error(),
				Severity: "CRITICAL",
				Trace:    trace,
			})

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Extract the characterIDs from the response.
		characterIDs := res2.Profiles.Response.Profile.Data.CharacterIds
		switch len(characterIDs) {
		case 0:
			characterIDs = []string{"", "", ""}
		case 1:
			characterIDs = []string{characterIDs[0], "", ""}
		case 2:
			characterIDs = []string{characterIDs[0], characterIDs[1], ""}
		case 3:
			characterIDs = []string{characterIDs[0], characterIDs[1], characterIDs[2]}
		}

		// Assemble the user.
		user := modelUser{
			Username:       username,
			Gamertag:       gamertag,
			MembershipType: membershipType,
			MembershipID:   membershipID,
			Characters: []modelCharacter{
				{ID: characterIDs[0]},
				{ID: characterIDs[1]},
				{ID: characterIDs[2]},
			},
		}

		// Create context for datastore.
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()

		// Create a datastore client.
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

		// Store the user into the database.
		key := datastore.IncompleteKey("User", nil)
		if _, err := c.Put(ctx, key, &user); err != nil {

			log.Println(Entry{
				Message:  err.Error(),
				Severity: "CRITICAL",
				Trace:    trace,
			})

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	// Execute the, create_user.html template.
	err := temps.ExecuteTemplate(w, "create_user.html", map[string]string{
		"CreateUserEndpoint": createUserEndpoint,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
