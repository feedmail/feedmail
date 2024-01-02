package user

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/feedmail/feedmail/app"
	M "github.com/feedmail/feedmail/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var user M.User

func New(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("user#new %v", r.URL)

	c.Respond(w, r, app.Tmpl{Handler: "user", Fn: "new", Layout: "layout/public"})

	return nil
}

func Create(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("user#create %v", r.URL)

	r.ParseForm()

	username := r.FormValue("username")
	if len(username) == 0 {
		return c.RespondErr(w, r, "shared", "username can't be blank")
	}

	email := r.FormValue("email")
	if len(email) == 0 {
		return c.RespondErr(w, r, "shared", "email can't be blank")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return c.RespondErr(w, r, "shared", "email is not valid")
	}

	password := r.FormValue("password")
	if len(password) == 0 {
		return c.RespondErr(w, r, "shared", "password can't be blank")
	}

	if len(password) < 8 {
		return c.RespondErr(w, r, "shared", "password must be 8 characters in length or longer")
	}

	res := c.DB.Client.Where("email = ?", strings.ToLower(email)).Find(&M.User{})
	if res.RowsAffected == 1 && strings.ToLower(email) == user.Email {
		//return app.RespondError(errors.New("email is already taken"))
		return c.RespondErr(w, r, "shared", "email is already taken")
	}

	userID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &M.User{
		ID:        userID,
		Username:  username,
		Email:     strings.ToLower(email),
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result := c.DB.Client.Create(&user)
	if result.Error != nil {
		//return app.RespondError(result.Error)
		return c.RespondErr(w, r, "shared", "can not create user")
	}

	// app.DB.Client.Model(&account).Association("Users").Append([]*S.User{user})
	// error handling?

	sessionID := uuid.New()
	session := &M.Session{
		ID:           sessionID,
		UserID:       user.ID,
		LastActivity: time.Now(),
	}
	result = c.DB.Client.Create(&session)
	if result.Error != nil {
		//return app.RespondError(result.Error)
		return c.RespondErr(w, r, "shared", "can not create session")
	}

	expiration := time.Now().Add(time.Minute * 30)

	csrfToken, err := c.CreateCsrfToken(sessionID.String())
	if err != nil {
		return c.RespondErr(w, r, "shared", "internal error")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    fmt.Sprintf("%s.%s", sessionID, csrfToken),
		Expires:  expiration,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func Update(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("user#update %v", r.URL)

	return nil
}

func Delete(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("user#delete %v", r.URL)

	return nil
}

func Profile(c *app.Config, w http.ResponseWriter, r *http.Request) error {
	log.Printf("user#profile %v", r.URL)

	c.Respond(w, r, app.Tmpl{Handler: "user", Fn: "profile"})

	return nil
}
