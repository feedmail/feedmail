package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Username  string `gorm:"type:VARCHAR(50);NOT NULL" json:"username"`
	Email     string `gorm:"type:VARCHAR(100);UNIQUE;NOT NULL" json:"email"`
	Password  []byte `json:"-"`
	Sessions  []Session
	Account   Account
}
