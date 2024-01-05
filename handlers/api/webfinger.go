package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	AP "github.com/feedmail/feedmail/activitypub"
	"github.com/feedmail/feedmail/app"
	M "github.com/feedmail/feedmail/models"
)

func Webfinger(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("api#webfinger %v", r.URL)

	resource := r.URL.Query().Get("resource")

	if len(resource) == 0 {
		return app.RespondStatus(w, http.StatusNotFound)
	}

	if !strings.HasPrefix(resource, "acct:") && !strings.HasSuffix(resource, *c.Domain) {
		return app.RespondStatus(w, http.StatusNotFound)
	}

	handle := strings.TrimPrefix(resource, "acct:")
	username := strings.TrimSuffix(handle, "@"+*c.Domain)

	var user M.User
	res := c.DB.Client.Where("username = ?", strings.ToLower(username)).Find(&user)
	if res.Error != nil || res.RowsAffected == 0 {
		return app.RespondStatus(w, http.StatusNotFound)
	}

	resp := AP.Webfinger{
		Subject: resource,
		Aliases: []string{
			fmt.Sprintf("https://%s/@%s", *c.Domain, user.Username),
			fmt.Sprintf("https://%s/users/%s", *c.Domain, user.Username),
		},
		Links: []AP.Link{
			{Rel: "http://webfinger.net/rel/profile-page", Type: "text/html", Href: fmt.Sprintf("https://%s/@%s", *c.Domain, user.Username)},
			{Rel: "self", Type: "application/activity+json", Href: fmt.Sprintf("https://%s/users/%s", *c.Domain, user.Username)},
			{Rel: "http://ostatus.org/schema/1.0/subscribe", Template: fmt.Sprintf("https://%s/authorize_interaction?uri={uri}", *c.Domain)},
			{Rel: "http://webfinger.net/rel/avatar", Type: "image/png", Href: "user.png"},
		},
	}

	return app.RespondJSON(w, http.StatusOK, resp)
}
