package repository

import "talkspace-api/modules/talkbot/entity"


type TalkbotCommandRepositoryInterface interface {
	SaveUserMessage(talkbot entity.Talkbot) error

}

type TalkbotQueryRepositoryInterface interface {
	GetUserMessages(userID string) ([]entity.Talkbot, error)
}