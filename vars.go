package destinyhome

import (
	"os"
	"text/template"
)

var (
	temps         *template.Template
	apiKey        = os.Getenv("BUNGIE_API_KEY")
	projectID     = os.Getenv("PROJECT_ID")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)
