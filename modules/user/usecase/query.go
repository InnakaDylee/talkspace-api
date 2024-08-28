package usecase

import (
	"errors"
	"talkspace-api/modules/user/entity"
	"talkspace-api/modules/user/repository"
	"talkspace-api/utils/constant"
)

type userQueryUsecase struct {
	userCommandRepository repository.UserCommandRepositoryInterface
	userQueryRepository   repository.UserQueryRepositoryInterface
}

func NewUserQueryUsecase(ucr repository.UserCommandRepositoryInterface, uqr repository.UserQueryRepositoryInterface) UserQueryUsecaseInterface {
	return &userQueryUsecase{
		userCommandRepository: ucr,
		userQueryRepository:   uqr,
	}
}

func (uqs *userQueryUsecase) GetUserByID(id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New(constant.ERROR_ID_INVALID)
	}

	userEntity, errGetID := uqs.userQueryRepository.GetUserByID(id)
	if errGetID != nil {
		return entity.User{}, errors.New(constant.ERROR_DATA_EMPTY)
	}
	
	return userEntity, nil
}

func (uqs *userQueryUsecase) GetRequestPremiumUsers() ([]entity.User, error) {
	users, err := uqs.userQueryRepository.GetRequestPremiumUsers()
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}
