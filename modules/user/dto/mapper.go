package dto

import "talkspace-api/modules/user/entity"

// Request
func UserRegisterRequestToUserEntity(request UserRegisterRequest) entity.User {
	return entity.User{
		Fullname:        request.Fullname,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func UserLoginRequestToUserEntity(request UserLoginRequest) entity.User {
	return entity.User{
		Email:    request.Email,
		Password: request.Password,
	}
}

func UserUpdateProfileRequestToUserEntity(request UserUpdateProfileRequest) entity.User {
	return entity.User{
		Fullname:       request.Fullname,
		Email:          request.Email,
		ProfilePicture: request.ProfilePicture,
		Gender:         request.Gender,
		Birthdate:      request.Birthdate,
		BloodType:      request.BloodType,
		Weight:         request.Weight,
		Height:         request.Height,
	}
}

func UserNewPasswordRequestToUserEntity(request UserNewPasswordRequest) entity.User {
	return entity.User{
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func UserUpdatePasswordRequestToUserEntity(request UserUpdatePasswordRequest) entity.User {
	return entity.User{
		Password:        request.Password,
		NewPassword:     request.NewPassword,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func UserSendOTPRequestToUserEntity(request UserSendOTPRequest) entity.User {
	return entity.User{
		Email: request.Email,
	}
}

func UserVerifyOTPRequestToUserEntity(request UserVerifyOTPRequest) entity.User {
	return entity.User{
		Email: request.Email,
		OTP:   request.OTP,
	}
}

// Response
func UserEntityToUserRegisterResponse(response entity.User) UserRegisterResponse {
	return UserRegisterResponse{
		ID:         response.ID,
		Fullname:   response.Fullname,
		Email:      response.Email,
	}
}

func UserEntityToUserLoginResponse(response entity.User, token string) UserLoginResponse {
	return UserLoginResponse{
		ID:         response.ID,
		Fullname:   response.Fullname,
		Email:      response.Email,
		Token:      token,
	}
}

func UserEntityToUserUpdateProfileResponse(response entity.User) UserUpdateProfileResponse {
	return UserUpdateProfileResponse{
		ID:             response.ID,
		Fullname:       response.Fullname,
		Email:          response.Email,
		ProfilePicture: response.ProfilePicture,
		Gender:         response.Gender,
		Birthdate:      response.Birthdate,
		BloodType:      response.BloodType,
		Weight:         response.Weight,
		Height:         response.Height,
	}
}

func UserEntityToUserProfileResponse(response entity.User) UserProfileResponse {
	return UserProfileResponse{
		ID:             response.ID,
		Fullname:       response.Fullname,
		Email:          response.Email,
		ProfilePicture: response.ProfilePicture,
		Gender:         response.Gender,
		Birthdate:      response.Birthdate,
		BloodType:      response.BloodType,
		Weight:         response.Weight,
		Height:         response.Height,
	}
}


func UserEntityToUserResponse(response entity.User) UserResponse {
	return UserResponse{
		ID:         response.ID,
		Fullname:   response.Fullname,
		Email:      response.Email,
	}
}

