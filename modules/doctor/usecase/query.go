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

func (dqs *doctorQueryUsecase) GetAllDoctors(status *bool, specialization string) ([]entity.Doctor, error) {
	doctors, err := dqs.doctorQueryRepository.GetAllDoctors(status, specialization)
	if err != nil {
		if err.Error() == constant.ERROR_DATA_EMPTY {
			return nil, errors.New(constant.ERROR_DATA_EMPTY)
		}
		return nil, err
	}

	return doctors, nil
}
