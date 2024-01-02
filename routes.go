package main

import (
	"github.com/feedmail/feedmail/app"
	"github.com/feedmail/feedmail/handlers/api/activitypub"
	"github.com/feedmail/feedmail/handlers/inbox"
	"github.com/feedmail/feedmail/handlers/sent"
	"github.com/feedmail/feedmail/handlers/session"
	"github.com/feedmail/feedmail/handlers/settings"
	"github.com/feedmail/feedmail/handlers/starred"
	"github.com/feedmail/feedmail/handlers/trash"
	"github.com/feedmail/feedmail/handlers/user"
)

func InitRoutes(c *app.Config) {
	// Public handlers
	c.Router.Handle("/sign-in", app.Public{Config: c, R: map[string]any{
		"GET":  session.New,
		"POST": session.Create,
	}})

	c.Router.Handle("/sign-up", app.Public{Config: c, R: map[string]any{
		"GET":  user.New,
		"POST": user.Create,
	}})

	// Auth handlers
	c.Router.Handle("/", app.Auth{Config: c, R: map[string]any{
		"GET": inbox.Index,
	}})

	c.Router.Handle("/starred", app.Auth{Config: c, R: map[string]any{
		"GET": starred.Index,
	}})

	c.Router.Handle("/sent", app.Auth{Config: c, R: map[string]any{
		"GET": sent.Index,
	}})

	c.Router.Handle("/trash", app.Auth{Config: c, R: map[string]any{
		"GET": trash.Index,
	}})

	c.Router.Handle("/sign-out", app.Auth{Config: c, R: map[string]any{
		"POST": session.Delete,
	}})

	c.Router.Handle("/settings", app.Auth{Config: c, R: map[string]any{
		"GET": settings.Index,
	}})

	c.Router.Handle("/user", app.Auth{Config: c, R: map[string]any{
		"UPDATE": user.Update,
		"DELETE": user.Delete,
	}})

	// API handlers
	c.Router.Handle("/.well-known/webfinger", app.API{Config: c, R: map[string]any{
		"GET": activitypub.Webfinger,
	}})

	c.Router.Handle("/users/{username}", app.API{Config: c, R: map[string]any{
		"GET": activitypub.Actor,
	}})

	c.Router.Handle("/users/{username}/inbox", app.API{Config: c, R: map[string]any{
		"GET": activitypub.Inbox,
	}})

	c.Router.Handle("/users/{username}/outbox", app.API{Config: c, R: map[string]any{
		"GET": activitypub.Outbox,
	}})
}
