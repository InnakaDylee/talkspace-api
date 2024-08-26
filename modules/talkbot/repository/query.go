package repository

import (
	"talkspace-api/modules/talkbot/entity"

	"gorm.io/gorm"
)

type talkbotQueryRepository struct {
	db *gorm.DB
}

func NewTalkbotQueryRepository(db *gorm.DB) TalkbotQueryRepositoryInterface {
	return &talkbotQueryRepository{
		db: db,
	}
}

func (tr *talkbotQueryRepository) GetUserMessages(userID string) ([]entity.Talkbot, error) {

	var talkbot []entity.Talkbot

	result := tr.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(3).Find(&talkbot)
	if result.Error != nil {
		return []entity.Talkbot{}, result.Error
	}

	return talkbot, nil
}
