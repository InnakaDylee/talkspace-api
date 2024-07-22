package model

import "time"

type User struct {
	ID             string `gorm:"primarykey"`
	Fullname       string `gorm:"not null"`
	Email          string `gorm:"not null"`
	Password       string `gorm:"not null"`
	ProfilePicture string
	Birthdate      string
	Gender         string `gorm:"type:gender"`
	BloodType      string `gorm:"type:blood_type"`
	Height         int
	Weight         int
	Role           string `gorm:"type:role;default:'user'"`
	OTP            string `gorm:"not null"`
	OTPExpiration  int64
	VerifyAccount  string
	IsVerified     bool `gorm:"not null;default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}
