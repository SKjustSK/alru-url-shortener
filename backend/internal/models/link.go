package models

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	LinkID    int64      `gorm:"primaryKey;autoIncrement" json:"link_id"`
	UserID    *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	LongURL   string     `gorm:"type:text;not null" json:"long_url"`
	ShortCode string     `gorm:"type:varchar(15);uniqueIndex;not null" json:"short_code"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `gorm:"default:now()" json:"created_at"`
}
