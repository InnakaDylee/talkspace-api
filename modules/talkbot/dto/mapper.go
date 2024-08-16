package dto

import "talkspace-api/modules/talkbot/entity"

// Request
func TalkBotRequestToTalkBotEntity(request TalkBotRequest) entity.TalkBot {
	return entity.TalkBot{
		Message: request.Message,
	}
}

// Response
func TalkBotEntityToTalkBotResponse(talkBot entity.TalkBot) TalkBotResponse {
	return TalkBotResponse{
		ID:        talkBot.ID,
		UserID:    talkBot.UserID,
		SessionID: talkBot.SessionID,
		Message:   talkBot.Message,
		Role:      talkBot.Role,
		CreatedAt: talkBot.CreatedAt, 
	}
}