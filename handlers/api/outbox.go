package api

import (
	"log"
	"net/http"

	AP "github.com/feedmail/feedmail/activitypub"
	"github.com/feedmail/feedmail/app"
)

func Outbox(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("api#outbox %v", r.URL)

	return app.RespondJSON(w, http.StatusOK, AP.Outbox{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
		},
		Id:           r.URL.Path,
		Type:         "OrderedCollection",
		TotalItems:   0,
		OrderedItems: []string{},
	})
}
