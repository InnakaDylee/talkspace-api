package model

import (
	"time"

	tm "talkspace-api/modules/transaction/model"
)

type Doctor struct {
	ID                string  `gorm:"primarykey"`
	Fullname          string  `gorm:"not null"`
	Email             string  `gorm:"not null"`
	Password          string  `gorm:"not null"`
	ProfilePicture    string  `gorm:"not null"`
	Status            bool    `gorm:"not null;default:true"`
	Gender            string  `gorm:"type:gender;default:NULL"`
	Specialization    string  `gorm:"not null"`
	YearsOfExperience string  `gorm:"not null"`
	Price             float64 `gorm:"not null"`
	LicenseNumber     string  `gorm:"not null"`
	Alumnus           string  `gorm:"not null"`
	About             string  `gorm:"not null"`
	Location          string  `gorm:"not null"`
	Role              string  `gorm:"type:role;default:'doctor'"`
	OTP               string  `gorm:"not null"`
	OTPExpiration     int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time       `gorm:"index"`
	Transaction       []tm.Transaction `gorm:"foreignKey:DoctorID"`
}
