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

	return app.RespondJSON(w, http.StatusOK, AP.Webfinger{
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
	})
}

func Nodefinger(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("api#nodefinger %v", r.URL)

	return app.RespondJSON(w, http.StatusOK, AP.Webfinger{
		Links: []AP.Link{
			{Rel: "http://nodeinfo.diaspora.software/ns/schema/2.0", Href: fmt.Sprintf("https://%s/%s", *c.Domain, "nodeinfo/2.0")},
		},
	})
}

func Nodeinfo(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("api#nodeinfo %v", r.URL)

	return app.RespondJSON(w, http.StatusOK, AP.Nodeinfo{
		Version: "2.0",
		Software: AP.Software{
			Name:    "feedmail",
			Version: "0.0.1",
		},
		Protocols: []string{"activitypub"},
		Services: AP.Services{
			Outbound: nil,
			Inbound:  nil,
		},
		Usage: AP.Usage{
			Users: AP.Users{
				Total: 0,
			},
			LocalPosts: 0,
		},
		OpenRegistrations: false,
	})
}
