package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
	UserID         uint
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
	SharedInboxUrl string
	ActorType      string
}
