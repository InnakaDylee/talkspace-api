package model

import "time"

type Consultation struct {
	ID            string `gorm:"primarykey"`
	TransactionID string `gorm:"not null"`
	SessionID     string `gorm:"not null"`
	UserID        string `gorm:"not null"`
	DoctorID      string `gorm:"not null"`
	Status 	  	  bool `gorm:"not null"`
	CreatedAt     time.Time
}

type Message struct {
	ID            string `gorm:"primarykey"`
	ConsultationID string `gorm:"not null"`
	ClientID	   string `gorm:"not null"`
	Message        string `gorm:"not null"`
	Role           string `gorm:"type:role;default:'user'"`
	CreatedAt      time.Time
}