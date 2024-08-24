package dto

type (
	UserRegisterResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
	}

	UserLoginResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
		Token      string `json:"token"`
	}

	UserUpdateProfileResponse struct {
		ID             string `json:"id"`
		Fullname       string `json:"fullname"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		Gender         string `json:"gender"`
		Birthdate      string `json:"birthdate"`
		BloodType      string `json:"blood_type"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
	}

	UserProfileResponse struct {
		ID             string `json:"id"`
		Fullname       string `json:"fullname"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
		Gender         string `json:"gender"`
		Birthdate      string `json:"birthdate"`
		BloodType      string `json:"blood_type"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
	}

	UserResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
	}

	UserListResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
		RequestPremium string `json:"request_premium"`
	}

)
