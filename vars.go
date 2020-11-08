package destinyhome

import (
	"os"
	"text/template"
)

var (
	temps              *template.Template
	apiKey             = os.Getenv("BUNGIE_API_KEY")
	projectID          = os.Getenv("PROJECT_ID")
	createUserEndpoint = os.Getenv("CREATE_USER_ENDPOINT")
)
