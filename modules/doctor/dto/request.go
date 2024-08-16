package dto

type (
	DoctorRegisterRequest struct {
		Fullname          string `json:"fullname" form:"fullname"`
		Email             string `json:"email" form:"email"`
		Password          string `json:"password" form:"password"`
		ProfilePicture    string `json:"profile_picture" form:"profile_picture"`
		Gender            string `json:"gender" form:"gender"`
		Specialization    string `json:"specialization" form:"specialization"`
		YearsOfExperience string `json:"years_of_experience" form:"years_of_experience"`
		LicenseNumber     string `json:"license_number" form:"license_number"`
		Alumnus           string `json:"alumnus" form:"alumnus"`
		About             string `json:"about" form:"about"`
		Location          string `json:"location" form:"location"`
	}

	DoctorLoginRequest struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	DoctorUpdateProfileRequest struct {
		Fullname          string `json:"fullname" form:"fullname"`
		Email             string `json:"email" form:"email"`
		ProfilePicture    string `json:"profile_picture" form:"profile_picture"`
		Gender            string `json:"gender" form:"gender"`
		Specialization    string `json:"specialization" form:"specialization"`
		LicenseNumber     string `json:"license_number" form:"license_number"`
		YearsOfExperience string    `json:"years_of_experience" form:"years_of_experience"`
		Alumnus           string `json:"alumnus" form:"alumnus"`
		About             string `json:"about" form:"about"`
		Location          string `json:"location" form:"location"`
	}

	DoctorUpdateStatusRequest struct {
		Status bool `json:"status" form:"status"`
	}

	DoctorNewPasswordRequest struct {
		Password        string `json:"password" form:"password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	DoctorUpdatePasswordRequest struct {
		Password        string `json:"password" form:"password"`
		NewPassword     string `json:"new_password" form:"new_password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	DoctorSendOTPRequest struct {
		Email string `json:"email" form:"email"`
	}

	DoctorVerifyOTPRequest struct {
		Email string `json:"email" form:"email"`
		OTP   string `json:"otp" form:"otp"`
	}
)
