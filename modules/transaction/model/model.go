package model

import (
	"time"
)

type Transaction struct {
	ID        string `gorm:"primarykey"`
	DoctorID  string `gorm:"foreignKey:DoctorID"`
	UserID    string `gorm:"foreignKey:UserID"`
	Status    bool   `gorm:"not null;default:false"`
	Amount    float64
	Method    string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
