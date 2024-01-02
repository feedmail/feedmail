package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;"`
	LastActivity time.Time
	UserID       uuid.UUID
}
