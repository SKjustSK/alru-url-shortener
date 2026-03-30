package models

import "time"

type Click struct {
	ClickID    int64     `gorm:"primaryKey;autoIncrement" json:"click_id"`
	LinkID     int64     `gorm:"index;not null" json:"link_id"`
	IPHash     string    `gorm:"type:varchar(64);not null" json:"-"`
	Country    string    `gorm:"type:varchar(2)" json:"country"`
	Referrer   string    `gorm:"type:text" json:"referrer"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	DeviceType string    `gorm:"type:varchar(50)" json:"device_type"`
	OS         string    `gorm:"type:varchar(50)" json:"os"`
	Browser    string    `gorm:"type:varchar(50)" json:"browser"`
	ClickedAt  time.Time `gorm:"autoCreateTime;index" json:"clicked_at"`
}
