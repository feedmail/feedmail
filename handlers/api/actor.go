package api

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"slices"
	"strings"

	AP "github.com/feedmail/feedmail/activitypub"
	"github.com/feedmail/feedmail/app"
	"github.com/feedmail/feedmail/handlers/user"
	M "github.com/feedmail/feedmail/models"
)

func Actor(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("api#actor %v", r.URL)

	headers := []string{
		"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
		"application/activity+json",
	}

	if !slices.Contains(headers, r.Header.Get("accept")) {
		return user.Profile(c, w, r) // redirect to user profile
	}

	username := r.PathValue("username")
	if len(username) == 0 {
		return app.RespondStatus(w, http.StatusNotFound)
	}

	var user M.User
	res := c.DB.Client.Where("username = ?", strings.ToLower(username)).Preload("Account").Find(&user)
	if res.Error != nil || res.RowsAffected == 0 {
		return app.RespondStatus(w, http.StatusNotFound)
	}

	id := fmt.Sprintf("https://%s/users/%s", *c.Domain, username)

	return app.RespondJSON(w, http.StatusOK, AP.Actor{
		Context: []string{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		Id:                id,
		Type:              "Person",
		Following:         path.Join(id, "following"),
		Followers:         path.Join(id, "followers"),
		Inbox:             path.Join(id, "inbox"),
		Outbox:            path.Join(id, "outbox"),
		PreferredUsername: username,
		Name:              username,
		Summary:           "",
		PublicKey: AP.PublicKey{
			Id:           path.Join(id, "#main-key"),
			Owner:        id,
			PublicKeyPem: user.Account.PublicKey,
		},
	})
}
