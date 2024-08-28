package usecase

import (
	"context"
	"talkspace-api/modules/talkbot/entity"

	"github.com/sashabaranov/go-openai"
)

type TalkbotQueryUsecaseInterface interface {
	GetTalkBotPrompt(userID string, talkbot entity.Talkbot, key string) (string, error)
	GetCompletionMessages(ctx context.Context, client *openai.Client, messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error)
}
