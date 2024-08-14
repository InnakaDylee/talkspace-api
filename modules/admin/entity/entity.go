package entity

import "time"

type Admin struct {
	ID              string
	Fullname        string
	Email           string
	Password        string
	NewPassword     string
	ConfirmPassword string
	Role            string
	OTP             string
	OTPExpiration   int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
