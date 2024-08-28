package usecase

import (
	"context"
	"os"
	"talkspace-api/modules/talkbot/entity"
	"talkspace-api/modules/talkbot/repository"

	"github.com/sashabaranov/go-openai"
)

type talkbotQueryUsecase struct {
	talkbotCommandRepository repository.TalkbotCommandRepositoryInterface
	talkbotQueryRepository   repository.TalkbotQueryRepositoryInterface
}

func NewTalkbotQueryUsecase(tcr repository.TalkbotCommandRepositoryInterface, tqr repository.TalkbotQueryRepositoryInterface) TalkbotQueryUsecaseInterface {
	return &talkbotQueryUsecase{
		talkbotCommandRepository: tcr,
		talkbotQueryRepository:   tqr,
	}
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

func (tqs *talkbotQueryUsecase) GetTalkBotPrompt(userID string, talkbot entity.Talkbot, key string) (string, error) {
	filePath := "utils/helper/prompt/talkbot-prompt-setup.txt"

	promptSetup, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	previousMessages, err := tqs.talkbotQueryRepository.GetUserMessages(userID)
	if err != nil {
		return "", err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: string(promptSetup),
		},
	}

	for _, msg := range previousMessages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "user",
			Content: msg.Message,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: talkbot.Message,
	})

	ctx := context.Background()
	client := openai.NewClient(key)
	model := openai.GPT3Dot5Turbo

	promptResponse, err := tqs.GetCompletionMessages(ctx, client, messages, model)
	if err != nil {
		return "", err
	}

	response := promptResponse.Choices[0].Message.Content
	userMessage := entity.Talkbot{
		UserID:  userID,
		Message: talkbot.Message,
	}

	if err := tqs.talkbotCommandRepository.SaveUserMessage(userMessage); err != nil {
		return "", err
	}

	return response, nil
}
