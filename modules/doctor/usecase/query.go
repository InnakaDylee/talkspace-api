package usecase

import (
	"errors"
	"talkspace-api/modules/doctor/entity"
	"talkspace-api/modules/doctor/repository"
	"talkspace-api/utils/constant"
)

type doctorQueryUsecase struct {
	doctorCommandRepository repository.DoctorCommandRepositoryInterface
	doctorQueryRepository   repository.DoctorQueryRepositoryInterface
}

func NewDoctorQueryUsecase(dcr repository.DoctorCommandRepositoryInterface, dqr repository.DoctorQueryRepositoryInterface) DoctorQueryUsecaseInterface {
	return &doctorQueryUsecase{
		doctorCommandRepository: dcr,
		doctorQueryRepository:   dqr,
	}
}

func (dqs *doctorQueryUsecase) GetDoctorByID(id string) (entity.Doctor, error) {
	if id == "" {
		return entity.Doctor{}, errors.New(constant.ERROR_ID_INVALID)
	}

	doctorEntity, errGetID := dqs.doctorQueryRepository.GetDoctorByID(id)
	if errGetID != nil {
		return entity.Doctor{}, errors.New(constant.ERROR_DATA_EMPTY)
	}
	
	return doctorEntity, nil
}

func (dqs *doctorQueryUsecase) GetAllDoctors(status *bool, specialization string, page, limit int) ([]entity.Doctor, int, error) {
	doctors, totalItems, err := dqs.doctorQueryRepository.GetAllDoctors(status, specialization, page, limit)
	if err != nil {
		if err.Error() == constant.ERROR_DATA_EMPTY {
			return nil, 0, errors.New(constant.ERROR_DATA_EMPTY)
		}
		return nil, 0, err
	}

	return doctors, totalItems, nil
}

