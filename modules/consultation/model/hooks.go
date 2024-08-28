package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (c *Consultation) BeforeCreate(tx *gorm.DB) (err error) {
	UUID := uuid.New()
	c.ID = UUID.String()

	// if c.Role == "" {
	// 	return errors.New("role is required")
	// }

	// validRoles := map[string]bool{"user": true, "doctor": true}
	// if !validRoles[c.Role] {
	// 	return errors.New("invalid role")
	// }

	return nil
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	UUID := uuid.New()
	m.ID = UUID.String()

	// if c.Role == "" {
	// 	return errors.New("role is required")
	// }

	// validRoles := map[string]bool{"user": true, "doctor": true}
	// if !validRoles[c.Role] {
	// 	return errors.New("invalid role")
	// }

	return nil
}