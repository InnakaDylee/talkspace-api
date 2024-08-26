package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (t *Talkbot) BeforeCreate(tx *gorm.DB) (err error) {
	UUID := uuid.New()
	t.ID = UUID.String()

	return nil
}
