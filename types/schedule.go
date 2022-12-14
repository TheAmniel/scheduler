package types

import "gorm.io/gorm"

type Schedule struct {
	ID        string `gorm:"primaryKey;not null" json:"id"`
	Status    string `gorm:"not null" json"-"`
	Content   string `gorm:"not null" json:"content,omitempty"`
	CreatedAt int64  `gorm:"type:bigint;not null" json:"created_at"`
	ExpiresAt int64  `gorm:"type:bigint;not null" json:"expires_at"`
}

func (sch *Schedule) BeforeCreate(tx *gorm.DB) error {
	if sch.Status != "Active" {
		sch.Status = "Active"
	}
	return nil
}
