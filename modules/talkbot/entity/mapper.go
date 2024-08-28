package entity

import "talkspace-api/modules/talkbot/model"

func TalkbotEntityToTalkbotModel(talkbotEntity Talkbot) model.Talkbot {
	return model.Talkbot{
		ID:        talkbotEntity.ID,
		UserID:    talkbotEntity.UserID,
		Message:   talkbotEntity.Message,
		CreatedAt: talkbotEntity.CreatedAt,
	}
}

func ListTalkbotEntityToTalkbotModel(talkbotEntities []Talkbot) []model.Talkbot {
	listTalkbotModel := []model.Talkbot{}
	for _, talkbot := range talkbotEntities {
		talkbotModel := TalkbotEntityToTalkbotModel(talkbot)
		listTalkbotModel = append(listTalkbotModel, talkbotModel)
	}
	return listTalkbotModel
}


func TalkbotModelToTalkbotEntity(talkbotModel model.Talkbot) Talkbot {
	return Talkbot{
		ID:        talkbotModel.ID,
		UserID:    talkbotModel.UserID,
		Message:   talkbotModel.Message,
		CreatedAt: talkbotModel.CreatedAt,
	}
}

func ListTalkbotModelToTalkbotEntity(talkbotModels []model.Talkbot) []Talkbot {
	listTalkbotEntity := []Talkbot{}
	for _, talkbot := range talkbotModels {
		talkbotEntity := TalkbotModelToTalkbotEntity(talkbot)
		listTalkbotEntity = append(listTalkbotEntity, talkbotEntity)
	}
	return listTalkbotEntity
}