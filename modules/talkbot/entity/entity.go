package entity

import (
	"time"
)

type Talkbot struct {
	ID        string
	UserID    string
	Message   string
	CreatedAt time.Time
}
