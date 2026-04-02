package models

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	LinkID    int64     `gorm:"primaryKey;autoIncrement" json:"link_id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	LongURL   string    `gorm:"type:text;not null" json:"long_url"`
	ShortCode string    `gorm:"type:varchar(15);not null;uniqueIndex:idx_shortcode_custom" json:"short_code"`
	IsCustom  bool      `gorm:"default:false;uniqueIndex:idx_shortcode_custom" json:"is_custom"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
