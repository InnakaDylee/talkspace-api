package model

import "time"

type Talkbot struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	Message   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
