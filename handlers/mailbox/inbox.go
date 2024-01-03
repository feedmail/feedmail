package mailbox

import (
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
)

func Inbox(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("mailbox#inbox %v", r.URL)

	return c.Respond(w, r, app.Tmpl{Handler: "inbox", Fn: "index"})
}
