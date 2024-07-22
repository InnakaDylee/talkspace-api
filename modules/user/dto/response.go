package dto

type (
	UserRegisterResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}

	UserLoginResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
		Token      string `json:"token"`
	}

	UserUpdateResponse struct {
		ID             string `json:"id"`
		Fullname       string `json:"fullname"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		Gender         string `json:"gender"`
		Birthdate      string `json:"birthdate"`
		BloodType      string `json:"blood_type"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
		IsVerified     bool   `json:"is_verified"`
	}

	UserResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}

)
