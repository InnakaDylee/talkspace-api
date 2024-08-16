package entity

import (
	"talkspace-api/modules/user/model"
	"time"
)

type TalkBot struct {
	ID        string
	UserID    string
	User      model.User
	SessionID string
	Message   string
	Role      string
	CreatedAt time.Time
}
