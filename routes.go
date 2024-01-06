package main

import (
	"github.com/feedmail/feedmail/app"
	"github.com/feedmail/feedmail/handlers/api"
	"github.com/feedmail/feedmail/handlers/mailbox"
	"github.com/feedmail/feedmail/handlers/session"
	"github.com/feedmail/feedmail/handlers/settings"
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
		"GET": user.Home,
	}})

	c.Router.Handle("/inbox", app.Auth{Config: c, R: map[string]any{
		"GET": mailbox.Inbox,
	}})

	c.Router.Handle("/starred", app.Auth{Config: c, R: map[string]any{
		"GET": mailbox.Starred,
	}})

	c.Router.Handle("/sent", app.Auth{Config: c, R: map[string]any{
		"GET": mailbox.Sent,
	}})

	c.Router.Handle("/trash", app.Auth{Config: c, R: map[string]any{
		"GET": mailbox.Trash,
	}})

	c.Router.Handle("/search", app.Auth{Config: c, R: map[string]any{
		"POST": mailbox.Search,
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
		"GET": api.Webfinger,
	}})

	c.Router.Handle("/users/{username}", app.API{Config: c, R: map[string]any{
		"GET": api.Actor,
	}})

	c.Router.Handle("/users/{username}/inbox", app.API{Config: c, R: map[string]any{
		"GET": api.Inbox,
	}})

	c.Router.Handle("/users/{username}/outbox", app.API{Config: c, R: map[string]any{
		"GET": api.Outbox,
	}})

	c.Router.Handle("/.well-known/nodeinfo", app.API{Config: c, R: map[string]any{
		"GET": api.Nodefinger,
	}})

	c.Router.Handle("/nodeinfo/2.0", app.API{Config: c, R: map[string]any{
		"GET": api.Nodeinfo,
	}})
}
