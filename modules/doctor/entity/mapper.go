package entity

import "talkspace-api/modules/doctor/model"

func DoctorEntityToDoctorModel(doctorEntity Doctor) model.Doctor {
	var gender *string
	if doctorEntity.Gender != "" {
		genderValue := doctorEntity.Gender
		gender = &genderValue
	}

	doctorModel := model.Doctor{
		ID:                doctorEntity.ID,
		Fullname:          doctorEntity.Fullname,
		Email:             doctorEntity.Email,
		Password:          doctorEntity.Password,
		ProfilePicture:    doctorEntity.ProfilePicture,
		Gender:            gender,
		Price:             doctorEntity.Price,
		Specialization:    doctorEntity.Specialization,
		LicenseNumber:     doctorEntity.LicenseNumber,
		YearsOfExperience: doctorEntity.YearsOfExperience,
		Alumnus:           doctorEntity.Alumnus,
		About:             doctorEntity.About,
		Location:          doctorEntity.Location,
		Status:            doctorEntity.Status,
		Role:              doctorEntity.Role,
		OTP:               doctorEntity.OTP,
		OTPExpiration:     doctorEntity.OTPExpiration,
		CreatedAt:         doctorEntity.CreatedAt,
		UpdatedAt:         doctorEntity.UpdatedAt,
		DeletedAt:         doctorEntity.DeletedAt,
	}
	return doctorModel
}

func ListDoctorEntityToDoctorModel(doctorEntities []Doctor) []model.Doctor {
	listDoctorModel := []model.Doctor{}
	for _, doctor := range doctorEntities {
		doctorModel := DoctorEntityToDoctorModel(doctor)
		listDoctorModel = append(listDoctorModel, doctorModel)
	}
	return listDoctorModel
}

func DoctorModelToDoctorEntity(doctorModel model.Doctor) Doctor {
	var gender string
	if doctorModel.Gender != nil {
		gender = *doctorModel.Gender
	}

	doctorEntity := Doctor{
		ID:                doctorModel.ID,
		Fullname:          doctorModel.Fullname,
		Email:             doctorModel.Email,
		Password:          doctorModel.Password,
		ProfilePicture:    doctorModel.ProfilePicture,
		Gender:            gender,
		Price:             doctorModel.Price,
		Specialization:    doctorModel.Specialization,
		LicenseNumber:     doctorModel.LicenseNumber,
		YearsOfExperience: doctorModel.YearsOfExperience,
		Alumnus:           doctorModel.Alumnus,
		About:             doctorModel.About,
		Location:          doctorModel.Location,
		Status:            doctorModel.Status,
		Role:              doctorModel.Role,
		OTP:               doctorModel.OTP,
		OTPExpiration:     doctorModel.OTPExpiration,
		CreatedAt:         doctorModel.CreatedAt,
		UpdatedAt:         doctorModel.UpdatedAt,
		DeletedAt:         doctorModel.DeletedAt,
	}
	return doctorEntity
}

func ListDoctorModelToDoctorEntity(doctorModels []model.Doctor) []Doctor {
	listDoctorEntity := []Doctor{}
	for _, doctor := range doctorModels {
		doctorEntity := DoctorModelToDoctorEntity(doctor)
		listDoctorEntity = append(listDoctorEntity, doctorEntity)
	}
	return listDoctorEntity
}
