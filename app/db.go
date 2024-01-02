package app

import (
	"net/http"
	"strings"

	M "github.com/feedmail/feedmail/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DB struct {
	Client *gorm.DB
}

func (db *DB) GetSessionID(r *http.Request) (bool, string) {
	sessionID, err := cookieSessionID(r)
	if err != nil {
		return false, ""
	}

	var session M.Session
	res := db.Client.Where("id = ?", sessionID).Where("user_id is not null").Find(&session)
	if res.Error != nil {
		return false, ""
	}
	return true, sessionID
}

func (db *DB) LoggedIn(r *http.Request) bool {
	sessionID, err := cookieSessionID(r)
	if err != nil {
		return false
	}

	var session M.Session
	res := db.Client.Where("id = ?", sessionID).Where("user_id is not null").Find(&session)
	return res.Error == nil
}

func (db *DB) GetCurrentUser(r *http.Request) (M.User, error) {
	var user M.User
	userID, err := db.GetCurrentUserID(r)
	if err != nil {
		return user, err
	}
	res := db.Client.Preload("Account").First(&user, userID)
	if res.RowsAffected == 0 {
		return user, res.Error
	}
	return user, nil
}

func (db *DB) GetCurrentUserID(r *http.Request) (uuid.UUID, error) {
	sessionID, err := cookieSessionID(r)
	if err != nil {
		return uuid.Nil, err
	}
	var session M.Session
	res := db.Client.Where("id = ?", sessionID).Where("user_id is not null").Find(&session)
	if res.Error != nil {
		return uuid.Nil, res.Error
	}
	return session.UserID, nil
}

func cookieSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	val := cookie.Value
	splitVal := strings.Split(val, ".")
	if len(splitVal) != 3 {
		return "", err
	}

	return splitVal[0], nil
}
