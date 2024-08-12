package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	UUID := uuid.New()
	u.ID = UUID.String()

	if u.Role == "" {
		u.Role = "user"
	}

	if u.Gender != nil {
		validGenders := map[string]bool{"male": true, "female": true}
		if !validGenders[*u.Gender] {
			return errors.New("invalid gender")
		}
	}

	if u.BloodType != nil {
		validBloodTypes := map[string]bool{"A": true, "B": true, "O": true, "AB": true}
		if !validBloodTypes[*u.BloodType] {
			return errors.New("invalid blood type")
		}
	}

	validRoles := map[string]bool{"user": true}
	if !validRoles[u.Role] {
		return errors.New("invalid role")
	}

	return nil
}

/*
CREATE TYPE gender AS ENUM ('male', 'female');
CREATE TYPE blood_type AS ENUM ('A', 'B', 'O', 'AB');
CREATE TYPE role AS ENUM ('user');
*/