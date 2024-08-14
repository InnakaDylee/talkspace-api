package model

import "time"

type Doctor struct {
	ID             string  `gorm:"primarykey"`
	Fullname       string  `gorm:"not null"`
	Email          string  `gorm:"not null"`
	Password       string  `gorm:"not null"`
	ProfilePicture string  `gorm:"not null"`
	Status         bool    `gorm:"not null;default:false"`
	Gender         *string `gorm:"type:gender;default:NULL"`
	Specialist     string  `gorm:"not null"`
	Experience     string  `gorm:"not null"`
	Price          float64 `gorm:"not null"`
	StrNumber      string  `gorm:"not null"`
	Alumnus        string  `gorm:"not null"`
	About          string  `gorm:"not null"`
	Location       string  `gorm:"not null"`
	Role           string  `gorm:"type:role;default:'doctor'"`
	OTP            string  `gorm:"not null"`
	OTPExpiration  int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}
