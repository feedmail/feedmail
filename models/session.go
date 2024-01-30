package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	LastActivity time.Time
	UserID       uuid.UUID `gorm:"type:uuid;"`
}
