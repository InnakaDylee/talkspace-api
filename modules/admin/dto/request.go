package dto

type (
	AdminRegisterRequest struct {
		Fullname        string `json:"fullname" form:"fullname"`
		Email           string `json:"email" form:"email"`
		Password        string `json:"password" form:"password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	AdminLoginRequest struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	AdminNewPasswordRequest struct {
		Password        string `json:"password" form:"password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	AdminUpdatePasswordRequest struct {
		Password        string `json:"password" form:"password"`
		NewPassword     string `json:"new_password" form:"new_password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	AdminSendOTPRequest struct {
		Email string `json:"email" form:"email"`
	}

	AdminVerifyOTPRequest struct {
		Email string `json:"email" form:"email"`
		OTP   string `json:"otp" form:"otp"`
	}
)
