package user

import (
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
)

func Home(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("user#home %v", r.URL)

	return c.Respond(w, r, app.Tmpl{Handler: "mailbox", Fn: "inbox"})
}
