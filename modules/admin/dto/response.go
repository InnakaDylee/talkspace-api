package dto

type (
	AdminRegisterResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
	}

	AdminLoginResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
		Token      string `json:"token"`
	}

	AdminResponse struct {
		ID         string `json:"id"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
	}
)
