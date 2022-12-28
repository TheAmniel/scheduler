package types

import "time"

type Schedule struct {
	ID        string        `gorm:"primaryKey;not null" json:"id"`
	Content   string        `gorm:"not null" json:"content,omitempty"`
	CreatedAt time.Duration `gorm:"type:bigint;not null" json:"created_at"`
	ExpiresAt time.Duration `gorm:"type:bigint;not null" json:"expires_at"`
}
