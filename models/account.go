package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
	UserID         uuid.UUID `gorm:"type:uuid;default:null;"`
	Username       string
	Domain         string
	PublicKey      string
	PrivateKey     string
	DisplayName    string
	Uri            string
	Url            string
	InboxUrl       string
	OutboxUrl      string
	FollowersUrl   string
	FollowingUrl   string
	SharedInboxUrl string
	ActorType      string
	IconUrl        string
	Summary        string
}
