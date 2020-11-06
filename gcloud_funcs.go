package destinyhome

import (
	"net/http"
)

// CreateUser is a Google Cloud Function for creating users.
// After it asks the user for their destiny info, it adds them
// to a Google Cloud Datastore, Database.
func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Execute the, create_user.html template.
	err := temps.ExecuteTemplate(w, "create_user.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MainWebHook(w http.ResponseWriter, r *http.Request) {}
