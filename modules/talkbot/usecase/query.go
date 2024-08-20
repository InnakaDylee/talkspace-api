package usecase

import (
	"context"
	"os"
	"talkspace-api/modules/talkbot/dto"

	"github.com/sashabaranov/go-openai"
)

type talkbotQueryUsecase struct{}

func NewTalkbotQueryUsecase() TalkbotQueryUsecaseInterface {
	return &talkbotQueryUsecase{}
}

func (tqs *talkbotQueryUsecase) GetCompletionMessages(ctx context.Context, client *openai.Client, messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error) {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	promptResponse, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	return promptResponse, err
}

func (tqs *talkbotQueryUsecase) GetTalkBotPrompt(request dto.TalkbotRequest, key string) (string, error) {
	filePath := "utils/helper/prompt/talkbot-prompt-setup.txt"

	promptSetup, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	client := openai.NewClient(key)
	model := openai.GPT3Dot5Turbo
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: string(promptSetup),
		},
		{
			Role:    "user",
			Content: request.Message,
		},
	}

	promptResponse, err := tqs.GetCompletionMessages(ctx, client, messages, model)
	if err != nil {
		return "", err
	}

	response := promptResponse.Choices[0].Message.Content
	return response, nil
}
