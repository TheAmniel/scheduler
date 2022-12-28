package types

type Schedule struct {
	ID        string `gorm:"primaryKey;not null" json:"id"`
	Content   string `gorm:"not null" json:"content,omitempty"`
	CreatedAt int64  `gorm:"type:bigint;not null" json:"created_at"`
	ExpiresAt int64  `gorm:"type:bigint;not null" json:"expires_at"`
}
