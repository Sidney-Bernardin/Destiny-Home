package destinyhome

import (
	"text/template"
)

func init() {

	// Parse tamplates.
	temps = template.Must(template.ParseGlob("templates/*.html"))
}
