package dto

import "talkspace-api/modules/admin/entity"

// Request
func AdminRegisterRequestToAdminEntity(request AdminRegisterRequest) entity.Admin {
	return entity.Admin{
		Fullname:        request.Fullname,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func AdminLoginRequestToAdminEntity(request AdminLoginRequest) entity.Admin {
	return entity.Admin{
		Email:    request.Email,
		Password: request.Password,
	}
}

func AdminNewPasswordRequestToAdminEntity(request AdminNewPasswordRequest) entity.Admin {
	return entity.Admin{
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func AdminUpdatePasswordRequestToAdminEntity(request AdminUpdatePasswordRequest) entity.Admin {
	return entity.Admin{
		Password:        request.Password,
		NewPassword:     request.NewPassword,
		ConfirmPassword: request.ConfirmPassword,
	}
}

func AdminSendOTPRequestToAdminEntity(request AdminSendOTPRequest) entity.Admin {
	return entity.Admin{
		Email: request.Email,
	}
}

func AdminVerifyOTPRequestToAdminEntity(request AdminVerifyOTPRequest) entity.Admin {
	return entity.Admin{
		Email: request.Email,
		OTP:   request.OTP,
	}
}

// Response
func AdminEntityToAdminRegisterResponse(response entity.Admin) AdminRegisterResponse {
	return AdminRegisterResponse{
		ID:       response.ID,
		Fullname: response.Fullname,
		Email:    response.Email,
	}
}

func AdminEntityToAdminLoginResponse(response entity.Admin, token string) AdminLoginResponse {
	return AdminLoginResponse{
		ID:       response.ID,
		Fullname: response.Fullname,
		Email:    response.Email,
		Token:    token,
	}
}

func AdminEntityToAdminResponse(response entity.Admin) AdminResponse {
	return AdminResponse{
		ID:       response.ID,
		Fullname: response.Fullname,
		Email:    response.Email,
	}
}
