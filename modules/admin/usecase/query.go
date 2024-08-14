package usecase

import (
	"errors"
	"talkspace-api/modules/admin/entity"      
	"talkspace-api/modules/admin/repository"  
	"talkspace-api/utils/constant"
)

type adminQueryUsecase struct {
	adminCommandRepository repository.AdminCommandRepositoryInterface 
	adminQueryRepository   repository.AdminQueryRepositoryInterface   
}

func NewAdminQueryUsecase(acr repository.AdminCommandRepositoryInterface, aqr repository.AdminQueryRepositoryInterface) AdminQueryUsecaseInterface {
	return &adminQueryUsecase{
		adminCommandRepository: acr,
		adminQueryRepository:   aqr,
	}
}

func (aqus *adminQueryUsecase) GetAdminByID(id string) (entity.Admin, error) { 
	if id == "" {
		return entity.Admin{}, errors.New(constant.ERROR_ID_INVALID)
	}

	adminEntity, errGetID := aqus.adminQueryRepository.GetAdminByID(id)
	if errGetID != nil {
		return entity.Admin{}, errors.New(constant.ERROR_DATA_EMPTY)
	}
	
	return adminEntity, nil
}
