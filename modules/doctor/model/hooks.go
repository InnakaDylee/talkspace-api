package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (d *Doctor) BeforeCreate(tx *gorm.DB) (err error) {
	UUID := uuid.New()
	d.ID = UUID.String()

	if d.Role == "" {
		d.Role = "doctor"
	}

	if d.Price == 0 {
		d.Price = 150000
	}

	validGenders := map[string]bool{"male": true, "female": true}
	if !validGenders[d.Gender] {
		return errors.New("invalid gender")
	}

	validRoles := map[string]bool{"doctor": true}
	if !validRoles[d.Role] {
		return errors.New("invalid role")
	}

	return nil
}

/*
CREATE TYPE gender AS ENUM ('Male', 'Female');
CREATE TYPE role AS ENUM ('doctor');
*/