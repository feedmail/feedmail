package activitypub

import (
	"log"
	"net/http"

	"github.com/feedmail/feedmail/app"
)

type outbox struct {
	Context      []string `json:"@context,omitempty"`
	Id           string   `json:"id,omitempty"`
	Type         string   `json:"type,omitempty"`
	TotalItems   int      `json:"totalItems"`
	OrderedItems []string `json:"orderedItems"`
}

func Outbox(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("activitypub#outbox %v", r.URL)

	return app.RespondJSON(w, http.StatusOK, outbox{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
		},
		Id:           r.URL.Path,
		Type:         "OrderedCollection",
		TotalItems:   0,
		OrderedItems: []string{},
	})
}
