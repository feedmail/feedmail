package starred

import (
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
)

func Index(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("starred#index %v", r.URL)

	return c.Respond(w, r, app.Tmpl{Handler: "starred", Fn: "index"})
}
