package repository

import (
	"talkspace-api/modules/talkbot/entity"

	"gorm.io/gorm"
)

type talkbotCommandRepository struct {
	db *gorm.DB
}

func NewTalkbotCommandRepository(db *gorm.DB) TalkbotCommandRepositoryInterface {
	return &talkbotCommandRepository{
		db: db,
	}
}

func (tr *talkbotCommandRepository) SaveUserMessage(talkbot entity.Talkbot) error {
	talkbotModel := entity.TalkbotEntityToTalkbotModel(talkbot)

	result := tr.db.Create(&talkbotModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
