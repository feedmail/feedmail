package session

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/feedmail/feedmail/app"
	M "github.com/feedmail/feedmail/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func New(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("session#new %v", r.URL)

	return c.Respond(w, r, app.Tmpl{Handler: "session", Fn: "new", Layout: "layout/public"})
}

func Create(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("session#create %v", r.URL)

	r.ParseForm()

	var user M.User
	res := c.DB.Client.Where("email = ?", strings.ToLower(r.FormValue("email"))).Find(&user)
	if res.Error != nil || res.RowsAffected == 0 {
		return c.RespondErr(w, r, "shared", "email or password are wrong")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.FormValue("password"))); err != nil {
		return c.RespondErr(w, r, "shared", "email or password are wrong")
	}

	sessionID := uuid.New()
	c.DB.Client.Create(&M.Session{
		ID:           sessionID,
		UserID:       user.ID,
		LastActivity: time.Now(),
	})

	expiration := time.Now().Add(time.Minute * 300)

	csrfToken, err := c.CreateCsrfToken(sessionID.String())
	if err != nil {
		return c.RespondErr(w, r, "shared", "internal error")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    fmt.Sprintf("%s.%s", sessionID, csrfToken),
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // HTTPS only.
	})
	// https://www.w3.org/TR/fetch-metadata/

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func Delete(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	ok, sessionID := c.DB.GetSessionID(r)
	if !ok {
		return app.RespondError(errors.New("you are not already logged in"))
	}

	var session M.Session
	c.DB.Client.Where("id = ?", sessionID).Delete(&session)

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
