package mailbox

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	AP "github.com/feedmail/feedmail/activitypub"
	"github.com/feedmail/feedmail/app"
	M "github.com/feedmail/feedmail/models"
)

type Actor struct {
	Name      string
	Icon      string
	Followers string
}

func Search(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("mailbox#search %v", r.URL)

	currentUser, err := c.DB.GetCurrentUser(r)
	if err != nil {
		return c.RespondErr(w, r, "shared", "you are not already logged in")
	}

	r.ParseForm()

	term := r.FormValue("term")
	if len(term) == 0 {
		return c.RespondErr(w, r, "shared", "username can't be blank")
	}

	// TODO: search for local users

	handle, _ := strings.CutPrefix(term, "@")
	split := strings.Split(strings.TrimSpace(handle), "@")
	if len(split) != 2 {
		log.Printf("wrong mail format")
	}
	remoteDomain := split[1]

	respWf, err := http.Get(fmt.Sprintf("https://%s/.well-known/webfinger?resource=acct:%s", remoteDomain, handle))
	if err != nil {
		log.Printf("error making http request: %s\n", err)
	}

	defer respWf.Body.Close()
	bodyWf, err := io.ReadAll(respWf.Body)
	if err != nil {
		log.Printf("can't read body")
	}

	var webfinger AP.Webfinger
	err = json.Unmarshal(bodyWf, &webfinger)
	if err != nil {
		log.Printf("can't parse body json")
	}

	var actorHref string
	for _, link := range webfinger.Links {
		if link.Rel == "self" && link.Type == "application/activity+json" {
			actorHref = link.Href
		}
	}

	client := http.Client{}
	reqActor, err := http.NewRequest("GET", actorHref, nil)
	if err != nil {
		log.Printf("can't create new request")
	}

	date := time.Now().UTC().Format(http.TimeFormat)
	headers := []AP.Pair{
		{K: "date", V: date},
		{K: "host", V: remoteDomain},
	}

	u, err := url.Parse(actorHref)
	if err != nil {
		log.Printf("can't parse user url")
	}

	keyId := fmt.Sprintf("https://%s/users/%s#main-key", *c.Domain, currentUser.Username)
	signature, err := AP.Sign(currentUser.Account.PrivateKey, keyId, "get "+u.Path, headers)
	if err != nil {
		log.Print("SignRequest", err)
	}

	reqActor.Header = http.Header{
		"Host":      []string{remoteDomain},
		"Date":      []string{date},
		"Accept":    {"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""},
		"Signature": []string{signature},
		//"Digest":    []string{digest},
	}
	log.Print(reqActor.Header)

	respActor, err := client.Do(reqActor)
	if err != nil {
		log.Print("can't query url", err)
	}

	defer respActor.Body.Close()
	bodyActor, err := io.ReadAll(respActor.Body)
	if err != nil {
		log.Print("can't read body", err)
	}

	log.Printf(">>> %s", bodyActor)

	remoteActor, err := AP.ParseActor(string(bodyActor))
	if err != nil {
		log.Print("can't parse body json", err)
	}

	uri := fmt.Sprintf("https://%s/users/%s", remoteDomain, handle)

	log.Print(uri)

	log.Printf(">>> %s", remoteActor)

	var account *M.Account
	getResult := c.DB.Client.Where("uri = ?", uri).Find(&account)
	if getResult.Error != nil || getResult.RowsAffected == 0 {
		// cache account
		account = &M.Account{
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Username:       remoteActor.PreferredUsername,
			Domain:         remoteDomain,
			PublicKey:      remoteActor.PublicKey.PublicKeyPem,
			DisplayName:    remoteActor.Name,
			Uri:            remoteActor.Id,
			Url:            remoteActor.Url,
			InboxUrl:       remoteActor.Inbox,
			OutboxUrl:      remoteActor.Outbox,
			FollowersUrl:   remoteActor.Followers,
			FollowingUrl:   remoteActor.Following,
			IconUrl:        remoteActor.Icon.Url,
			SharedInboxUrl: remoteActor.Endpoints.SharedInbox,
			Summary:        remoteActor.Summary,
			ActorType:      "Person",
		}
		accountResult := c.DB.Client.Create(&account)
		if accountResult.Error != nil {
			log.Print(accountResult.Error)
			return c.RespondErr(w, r, "shared", "can not create session")
		}
	} else {
		// update cached account
		account.PublicKey = remoteActor.PublicKey.PublicKeyPem
		account.DisplayName = remoteActor.Name
		account.Username = remoteActor.PreferredUsername
		account.Url = remoteActor.Url
		account.UpdatedAt = time.Now()
		updateResult := c.DB.Client.Save(&account)
		if updateResult.Error != nil {
			log.Printf("can't update cached remote actor")
		}
	}

	return c.Respond(w, r, app.Tmpl{Handler: "mailbox", Fn: "search", Action: "replace", Target: "main", Data: remoteActor})
}
