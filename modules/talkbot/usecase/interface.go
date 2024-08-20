package usecase

import (
	"context"
	"talkspace-api/modules/talkbot/dto"

	"github.com/sashabaranov/go-openai"
)

type TalkbotQueryUsecaseInterface interface {
	GetTalkBotPrompt(request dto.TalkbotRequest, key string) (string, error)
	GetCompletionMessages(ctx context.Context, client *openai.Client, messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error)
}
