package model

import (
	"talkspace-api/modules/user/model"
	"time"
)

type Talkbot struct {
	ID        string     `gorm:"primarykey"`
	UserID    string     `gorm:"index"`
	User      model.User `gorm:"foreignKey:UserID"`
	SessionID string     `gorm:"not null"`
	Message   string     `gorm:"not null"`
	Role      string     `gorm:"not null"`
	CreatedAt time.Time
}
