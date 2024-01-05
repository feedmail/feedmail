package mailbox

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	AP "github.com/feedmail/feedmail/activitypub"
	"github.com/feedmail/feedmail/app"
)

type Actor struct {
	Name      string
	Icon      string
	Followers string
}

func Search(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("mailbox#search %v", r.URL)

	r.ParseForm()

	term := r.FormValue("term")
	if len(term) == 0 {
		return c.RespondErr(w, r, "shared", "username can't be blank")
	}

	//var data string

	handle, _ := strings.CutPrefix(term, "@")
	split := strings.Split(handle, "@")
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

	log.Print(actorHref)

	client := http.Client{}
	reqActor, err := http.NewRequest("GET", actorHref, nil)
	if err != nil {
		log.Printf("can't create new request")
	}

	reqActor.Header = http.Header{
		"Accept": {"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""},
		//"Authorization": {"Bearer Token"},
	}

	respActor, err := client.Do(reqActor)
	if err != nil {
		log.Printf("can't query url")
	}

	defer respActor.Body.Close()
	bodyActor, err := io.ReadAll(respActor.Body)
	if err != nil {
		log.Printf("can't read body")
	}

	remoteActor := struct{ AP.Actor }{}

	err = json.Unmarshal(bodyActor, &remoteActor)
	if err != nil {
		log.Printf("can't parse body json")
	}

	// id := fmt.Sprintf("https://%s/users/%s", remoteDomain, handle)
	// accountID := uuid.New()
	// account := &M.Account{
	// 	ID:             accountID,
	// 	CreatedAt:      time.Now(),
	// 	UpdatedAt:      time.Now(),
	// 	Username:       remoteActor.PreferredUsername,
	// 	Domain:         remoteDomain,
	// 	PublicKey:      remoteActor.PublicKey.PublicKeyPem,
	// 	DisplayName:    remoteActor.Name,
	// 	Uri:            "",
	// 	Url:            remoteActor.Url,
	// 	InboxUrl:       path.Join(id, "inbox"),
	// 	OutboxUrl:      path.Join(id, "outbox"),
	// 	FollowersUrl:   path.Join(id, "followers"),
	// 	SharedInboxUrl: fmt.Sprintf("https://%s/inbox", remoteDomain),
	// 	ActorType:      "Person",
	// }
	// accountResult := c.DB.Client.Create(&account)
	// if accountResult.Error != nil {
	// 	log.Print(accountResult.Error)
	// 	//rollback user creation
	// 	return c.RespondErr(w, r, "shared", "can not create session")
	// }

	//log.Print(fmt.Sprintf("%s", actor["id"]))

	return c.Respond(w, r, app.Tmpl{Handler: "mailbox", Fn: "search", Action: "replace", Target: "main", Data: remoteActor})
}
