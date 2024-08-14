package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	UUID := uuid.New()
	a.ID = UUID.String()

	if a.Role == "" {
		a.Role = "admin"
	}

	validRoles := map[string]bool{"admin": true}
	if !validRoles[a.Role] {
		return errors.New("invalid role")
	}

	return nil
}

/*
CREATE TYPE role AS ENUM ('admin');
*/
