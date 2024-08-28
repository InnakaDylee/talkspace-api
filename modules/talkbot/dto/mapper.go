package dto

import "talkspace-api/modules/talkbot/entity"

// Request
func TalkbotRequestToTalkbotEntity(request TalkbotRequest) entity.Talkbot {
	return entity.Talkbot{
		Message: request.Message,
	}
}

// Response
func TalkbotEntityToTalkbotResponse(entity entity.Talkbot) TalkbotResponse {
	return TalkbotResponse{
		Message: entity.Message,
	}
}
