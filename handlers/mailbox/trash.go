package mailbox

import (
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
)

func Trash(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("mailbox#trash %v", r.URL)

	return c.Respond(w, r, app.Tmpl{Handler: "mailbox", Fn: "trash"})
}
