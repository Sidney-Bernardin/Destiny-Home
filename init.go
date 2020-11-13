package destinyhome

import (
	"log"
	"text/template"
)

func init() {

	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)

	// Parse tamplates.
	temps = template.Must(template.ParseGlob("templates/*.html"))
}
