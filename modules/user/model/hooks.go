package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	UUID := uuid.New()
	u.ID = UUID.String()

	validGenders := map[string]bool{"male": true, "female": true}
	if !validGenders[u.Gender] {
		return errors.New("invalid gender")
	}

	validBloodTypes := map[string]bool{"A": true, "B": true, "O": true, "AB": true}
	if !validBloodTypes[u.BloodType] {
		return errors.New("invalid blood type")
	}

	validRoles := map[string]bool{"user": true}
	if !validRoles[u.Role] {
		return errors.New("invalid role")
	}

	return nil
}
