package model

import "time"

type User struct {
	ID             string `gorm:"primarykey"`
	Fullname       string `gorm:"not null"`
	Email          string `gorm:"not null"`
	Password       string `gorm:"not null"`
	ProfilePicture string
	Birthdate      string
	Gender         *string `gorm:"type:gender;default:NULL"`
	BloodType      *string `gorm:"type:blood_type;default:NULL"`
	Height         int
	Weight         int
	Role           string `gorm:"type:role;default:'user'"`
	RequestPremium string
	PremiumExpired time.Time
	OTP            string `gorm:"not null"`
	OTPExpiration  int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}
