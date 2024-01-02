package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Public struct {
	*Config
	R map[string]any
}

func (p Public) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("PublicHandler %v", r.URL)

	if _, ok := p.R[r.Method]; !ok {
		log.Print("http method not allowed")
		HandleError(w, fmt.Errorf("http method not allowed"))
		return
	}

	err := p.R[r.Method].(func(e *Config, w http.ResponseWriter, r *http.Request) error)(p.Config, w, r)
	if err != nil {
		log.Print(err)
		HandleError(w, err)
		return
	}
}

type Auth struct {
	*Config
	R map[string]any
}

func (a Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("AuthHandler %v\n%v", r, a.R)

	if _, ok := a.R[r.Method]; !ok {
		log.Print("http method not allowed")
		HandleError(w, fmt.Errorf("http method not allowed"))
		return
	}

	ok := a.Config.DB.LoggedIn(r)
	if !ok {
		log.Print("Unauthorized")
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	// validate csrf token on form post
	if r.Method == http.MethodPost {
		if !a.Config.ValidCsrfToken(r) {
			log.Print("csrf token not valid")
			HandleError(w, fmt.Errorf("csrf token not valid"))
			return
		}
	}

	token, err := a.Config.GetCsrfToken(r)
	if len(token) == 0 || err != nil {
		log.Print("csrf token missing")
		HandleError(w, fmt.Errorf("csrf token missing"))
		return
	}

	ctx := context.WithValue(r.Context(), "csrf", token)
	newReq := r.WithContext(ctx)

	log.Printf("token %s", token)

	err = a.R[r.Method].(func(e *Config, w http.ResponseWriter, r *http.Request) error)(a.Config, w, newReq)
	if err != nil {
		log.Print(err)
		HandleError(w, err)
		return
	}
}

type API struct {
	*Config
	R map[string]any
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("APIHandler %v", r.URL)

	if _, ok := api.R[r.Method]; !ok {
		log.Print("http method not allowed")
		HandleError(w, fmt.Errorf("http method not allowed"))
		return
	}

	err := api.R[r.Method].(func(e *Config, w http.ResponseWriter, r *http.Request) error)(api.Config, w, r)
	if err != nil {
		log.Print(err)
		HandleError(w, err)
		return
	}
}
