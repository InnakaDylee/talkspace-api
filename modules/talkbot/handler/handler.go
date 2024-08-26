package handler

import (
	"net/http"
	"talkspace-api/app/configs"
	"talkspace-api/middlewares"
	"talkspace-api/modules/talkbot/dto"
	"talkspace-api/modules/talkbot/usecase"
	"talkspace-api/utils/constant"
	"talkspace-api/utils/responses"

	"github.com/labstack/echo/v4"
)

type talkbotHandler struct {
	talkbotQueryUsecase   usecase.TalkbotQueryUsecaseInterface
}

func NewTalkbotHandler(tqu usecase.TalkbotQueryUsecaseInterface) *talkbotHandler {
	return &talkbotHandler{
		talkbotQueryUsecase:   tqu,
	}
}

// Command
func (th *talkbotHandler) CreateTalkBotMessage(c echo.Context) error {
	talkbotRequest := dto.TalkbotRequest{}

	errBind := c.Bind(&talkbotRequest)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errBind.Error()))
	}

	userID, role, errExtractToken := middlewares.ExtractToken(c)
	if role != constant.USER {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}
	if errExtractToken != nil {
		return c.JSON(http.StatusUnauthorized, responses.ErrorResponse(errExtractToken.Error()))
	}

	config, err := configs.LoadConfig()
	if err != nil {
		return err
	}

	talkbotEntity :=  dto.TalkbotRequestToTalkbotEntity(talkbotRequest)

	promptResponse, errGetPrompt := th.talkbotQueryUsecase.GetTalkBotPrompt(userID, talkbotEntity, config.OPENAI.OPENAI_API_KEY)
	if errGetPrompt != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse(errGetPrompt.Error()))
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse(constant.SUCCESS_RETRIEVED, promptResponse))
}
