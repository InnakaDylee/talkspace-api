package dto

import "time"

type TalkBotResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	Message   string    `json:"message"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

