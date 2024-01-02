package activitypub

import (
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
)

func Inbox(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("activitypub#inbox %v", r.URL)

	return app.RespondStatus(w, http.StatusOK)
}
