package dto

import "talkspace-api/modules/doctor/entity"

// Request
func DoctorRegisterRequestToDoctorEntity(request DoctorRegisterRequest) entity.Doctor {
	return entity.Doctor{
		Fullname:          request.Fullname,
		Email:             request.Email,
		Password:          request.Password,
		ProfilePicture:    request.ProfilePicture,
		Gender:            request.Gender,
		Specialization:    request.Specialization,
		YearsOfExperience: request.YearsOfExperience,
		LicenseNumber:     request.LicenseNumber,
		Alumnus:           request.Alumnus,
		About:             request.About,
		Location:          request.Location,
	}
}

func DoctorLoginRequestToDoctorEntity(request DoctorLoginRequest) entity.Doctor {
	return entity.Doctor{
		Email:    request.Email,
		Password: request.Password,
	}
}

func DoctorUpdateProfileRequestToDoctorEntity(request DoctorUpdateProfileRequest) entity.Doctor {
	return entity.Doctor{
		Fullname:          request.Fullname,
		Email:             request.Email,
		ProfilePicture:    request.ProfilePicture,
		Gender:            request.Gender,
		Specialization:    request.Specialization,
		LicenseNumber:     request.LicenseNumber,
		YearsOfExperience: request.YearsOfExperience,
		Alumnus:           request.Alumnus,
		About:             request.About,
		Location:          request.Location,
	}
}

func DoctorUpdateStatusRequestToDoctorEntity(request DoctorUpdateStatusRequest) entity.Doctor {
	return entity.Doctor{
		Status: request.Status,
	}
}

func DoctorNewPasswordRequestToDoctorEntity(request DoctorNewPasswordRequest) entity.Doctor {
	return entity.Doctor{
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func DoctorUpdatePasswordRequestToDoctorEntity(request DoctorUpdatePasswordRequest) entity.Doctor {
	return entity.Doctor{
		Password:        request.Password,
		NewPassword:     request.NewPassword,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func DoctorSendOTPRequestToDoctorEntity(request DoctorSendOTPRequest) entity.Doctor {
	return entity.Doctor{
		Email: request.Email,
	}
}

func DoctorVerifyOTPRequestToDoctorEntity(request DoctorVerifyOTPRequest) entity.Doctor {
	return entity.Doctor{
		Email: request.Email,
		OTP:   request.OTP,
	}
}

// Response
func DoctorEntityToDoctorRegisterResponse(response entity.Doctor) DoctorRegisterResponse {
	return DoctorRegisterResponse{
		ID:                response.ID,
		Fullname:          response.Fullname,
		Email:             response.Email,
		ProfilePicture:    response.ProfilePicture,
		Status:            response.Status,
		Gender:            response.Gender,
		Specialization:    response.Specialization,
		YearsOfExperience: response.YearsOfExperience,
		LicenseNumber:     response.LicenseNumber,
		Alumnus:           response.Alumnus,
		About:             response.About,
		Location:          response.Location,
	}
}

func DoctorEntityToDoctorLoginResponse(entity entity.Doctor, token string) DoctorLoginResponse {
	return DoctorLoginResponse{
		ID:       entity.ID,
		Fullname: entity.Fullname,
		Email:    entity.Email,
		Token:    token,
	}
}

func DoctorEntityToDoctorUpdateProfileResponse(entity entity.Doctor) DoctorUpdateProfileResponse {
	return DoctorUpdateProfileResponse{
		ID:                entity.ID,
		Fullname:          entity.Fullname,
		Email:             entity.Email,
		ProfilePicture:    entity.ProfilePicture,
		Gender:            entity.Gender,
		Specialization:    entity.Specialization,
		LicenseNumber:     entity.LicenseNumber,
		YearsOfExperience: entity.YearsOfExperience,
		Alumnus:           entity.Alumnus,
		About:             entity.About,
		Location:          entity.Location,
	}
}

func DoctorEntityToDoctorProfileResponse(entity entity.Doctor) DoctorProfileResponse {
	return DoctorProfileResponse{
		ID:                entity.ID,
		Status:            entity.Status,
		Fullname:          entity.Fullname,
		Email:             entity.Email,
		ProfilePicture:    entity.ProfilePicture,
		Gender:            entity.Gender,
		Specialization:    entity.Specialization,
		LicenseNumber:     entity.LicenseNumber,
		YearsOfExperience: entity.YearsOfExperience,
		Alumnus:           entity.Alumnus,
		About:             entity.About,
		Location:          entity.Location,
	}
}

func DoctorEntityToDoctorUpdateStatusResponse(entity entity.Doctor) DoctorUpdateStatusResponse {
	return DoctorUpdateStatusResponse{
		ID:     entity.ID,
		Status: entity.Status,
	}
}

func DoctorEntityToDoctorResponse(entity entity.Doctor) DoctorResponse {
	return DoctorResponse{
		ID:       entity.ID,
		Fullname: entity.Fullname,
		Email:    entity.Email,
	}
}


