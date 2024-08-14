package model

import "time"

type Admin struct {
	ID            string `gorm:"primarykey"`
	Fullname      string `gorm:"not null"`
	Email         string `gorm:"not null"`
	Password      string `gorm:"not null"`
	Role          string `gorm:"type:role;default:'admin'"`
	OTP           string `gorm:"not null"`
	OTPExpiration int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `gorm:"index"`
}
