package user

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"path"
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
	userResult := c.DB.Client.Create(&user)
	if userResult.Error != nil {
		log.Print(userResult.Error)
		return c.RespondErr(w, r, "shared", "can not create user")
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Print("can not generate rsa key")
		return c.RespondErr(w, r, "shared", "can not create user")
	}

	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privKeyBytes,
		},
	)

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		log.Print("can not marshal public key")
		return c.RespondErr(w, r, "shared", "can not create user")
	}
	pubKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	)

	id := fmt.Sprintf("https://%s/users/%s", *c.Domain, username)
	accountID := uuid.New()
	account := &M.Account{
		ID:             accountID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		UserID:         user.ID,
		Username:       username,
		Domain:         *c.Domain,
		PublicKey:      string(pubKeyPem),
		PrivateKey:     string(privKeyPem),
		DisplayName:    username,
		Uri:            "",
		Url:            fmt.Sprintf("https://%s/@%s", *c.Domain, username),
		InboxUrl:       path.Join(id, "inbox"),
		OutboxUrl:      path.Join(id, "outbox"),
		FollowersUrl:   path.Join(id, "followers"),
		SharedInboxUrl: fmt.Sprintf("https://%s/inbox", *c.Domain),
		ActorType:      "Person",
	}
	accountResult := c.DB.Client.Create(&account)
	if accountResult.Error != nil {
		log.Print(accountResult.Error)
		//rollback user creation
		return c.RespondErr(w, r, "shared", "can not create session")
	}

	sessionID := uuid.New()
	session := &M.Session{
		ID:           sessionID,
		UserID:       user.ID,
		LastActivity: time.Now(),
	}
	sessionResult := c.DB.Client.Create(&session)
	if sessionResult.Error != nil {
		log.Print(sessionResult.Error)
		//rollback user and account creation
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
