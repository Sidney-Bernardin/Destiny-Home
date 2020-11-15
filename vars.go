package destinyhome

import (
	"os"
	"text/template"
)

var (
	apiKey        = os.Getenv("BUNGIE_API_KEY")
	projectID     = os.Getenv("PROJECT_ID")
	adminEndpoint = os.Getenv("ADMIN_ENDPOINT")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
	temps         *template.Template
)
