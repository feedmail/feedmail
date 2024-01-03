package mailbox

import (
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
)

func Starred(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("mailbox#starred %v", r.URL)

	return c.Respond(w, r, app.Tmpl{Handler: "starred", Fn: "index"})
}
