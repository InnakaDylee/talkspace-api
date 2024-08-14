package entity

import "time"

type Doctor struct {
	ID                string
	Fullname          string
	Email             string
	Password          string
	NewPassword       string
	ConfirmPassword   string
	ProfilePicture    string
	Gender            string
	Price             float64
	Specialization    string
	LicenseNumber     string
	YearsOfExperience string
	Alumnus           string
	About             string
	Location          string
	Status            bool
	Role              string
	OTP               string
	OTPExpiration     int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
