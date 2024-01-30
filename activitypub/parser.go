package activitypub

import (
	"encoding/json"
	"errors"
	"fmt"
)

func ParseActor(in string) (*Actor, error) {
	var a interface{}

	if err := json.Unmarshal([]byte(in), &a); err != nil {
		return nil, errors.New("can't unmarshal actor string")
	}

	actor := Actor{}

	if m, ok := a.(map[string]interface{}); ok {
		if preferredUsername, ok := m["preferredUsername"].(string); ok {
			actor.PreferredUsername = preferredUsername
		} else {
			fmt.Printf("can't parse preferredUsername %+v\n", m["preferredUsername"])
		}

		pk, err := ParseActorPublicKey(m["publicKey"])
		if err != nil {
			fmt.Printf("can't parse actor public key %v\n", err)
		} else {
			actor.PublicKey = *pk
		}

		icon, err := ParseActorIcon(m["icon"])
		if err != nil {
			fmt.Printf("can't parse actor icon %v\n", err)
		} else {
			actor.Icon = *icon
		}

		ep, err := ParseActorEndpoints(m["endpoints"])
		if err != nil {
			fmt.Printf("can't parse actor endpoints %v\n", err)
		} else {
			actor.Endpoints = *ep
		}

		if id, ok := m["id"].(string); ok {
			actor.Id = id
		} else {
			fmt.Printf("can't parse id %+v\n", m["id"])
		}

		if name, ok := m["name"].(string); ok {
			actor.Name = name
		} else {
			fmt.Printf("can't parse name %+v\n", m["name"])
		}

		if summary, ok := m["summary"].(string); ok {
			actor.Summary = summary
		} else {
			fmt.Printf("can't parse summary %+v\n", m["summary"])
		}

		if url, ok := m["url"].(string); ok {
			actor.Url = url
		} else {
			fmt.Printf("can't parse url %+v\n", m["url"])
		}

		if inbox, ok := m["inbox"].(string); ok {
			actor.Inbox = inbox
		} else {
			fmt.Printf("can't parse inbox %+v\n", m["inbox"])
		}

		if outbox, ok := m["outbox"].(string); ok {
			actor.Outbox = outbox
		} else {
			fmt.Printf("can't parse outbox %+v\n", m["outbox"])
		}

		if following, ok := m["following"].(string); ok {
			actor.Following = following
		} else {
			fmt.Printf("can't parse following %+v\n", m["following"])
		}

		if followers, ok := m["followers"].(string); ok {
			actor.Followers = followers
		} else {
			fmt.Printf("can't parse followers %+v\n", m["followers"])
		}

	} else {
		return nil, errors.New("can't unmarshal actor string with unexpected type")
	}

	return &actor, nil
}

func ParseActorPublicKey(in interface{}) (*PublicKey, error) {
	pk := PublicKey{}

	if m, ok := in.(map[string]interface{}); ok {
		if id, ok := m["id"].(string); ok {
			pk.Id = id
		} else {
			return nil, errors.New("can't parse actor public key id")
		}

		if owner, ok := m["owner"].(string); ok {
			pk.Owner = owner
		} else {
			return nil, errors.New("can't parse actor public key owner")
		}

		if publicKeyPem, ok := m["publicKeyPem"].(string); ok {
			pk.PublicKeyPem = publicKeyPem
		} else {
			return nil, errors.New("can't parse actor public key pem")
		}
	}

	return &pk, nil
}

func ParseActorIcon(in interface{}) (*Icon, error) {
	icon := Icon{}

	if m, ok := in.(map[string]interface{}); ok {
		if t, ok := m["type"].(string); ok {
			icon.Type = t
		} else {
			return nil, errors.New("can't parse actor icon id")
		}

		if mediaType, ok := m["mediaType"].(string); ok {
			icon.MediaType = mediaType
		} else {
			return nil, errors.New("can't parse actor icon mediaType")
		}

		if url, ok := m["url"].(string); ok {
			icon.Url = url
		} else {
			return nil, errors.New("can't parse actor icon url")
		}
	}

	return &icon, nil
}

func ParseActorEndpoints(in interface{}) (*Endpoints, error) {
	ep := Endpoints{}

	if m, ok := in.(map[string]interface{}); ok {
		if sharedInbox, ok := m["sharedInbox"].(string); ok {
			ep.SharedInbox = sharedInbox
		} else {
			return nil, errors.New("can't parse actor endpoints sharedInbox")
		}
	}

	return &ep, nil
}
