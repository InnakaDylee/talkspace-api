package dto

type (
	DoctorRegisterResponse struct {
		ID                string `json:"id"`
		Fullname          string `json:"fullname"`
		Email             string `json:"email"`
		ProfilePicture    string `json:"profile_picture"`
		Status            bool   `json:"status"`
		Gender            string `json:"gender"`
		Specialization    string `json:"specialization"`
		YearsOfExperience string `json:"years_of_experience"`
		LicenseNumber     string `json:"license_number"`
		Alumnus           string `json:"alumnus"`
		About             string `json:"about"`
		Location          string `json:"location"`
	}

	DoctorLoginResponse struct {
		ID       string `json:"id"`
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	}

	DoctorUpdateProfileResponse struct {
		ID                string `json:"id"`
		Fullname          string `json:"fullname"`
		Email             string `json:"email"`
		ProfilePicture    string `json:"profile_picture"`
		Gender            string `json:"gender"`
		Specialization    string `json:"specialization"`
		LicenseNumber     string `json:"license_number"`
		YearsOfExperience string `json:"years_of_experience"`
		Alumnus           string `json:"alumnus"`
		About             string `json:"about"`
		Location          string `json:"location"`
	}

	DoctorProfileResponse struct {
		ID                string `json:"id"`
		Status            bool   `json:"status"`
		Fullname          string `json:"fullname"`
		Email             string `json:"email"`
		ProfilePicture    string `json:"profile_picture"`
		Gender            string `json:"gender"`
		Specialization    string `json:"specialization"`
		LicenseNumber     string `json:"license_number"`
		YearsOfExperience string `json:"years_of_experience"`
		Alumnus           string `json:"alumnus"`
		About             string `json:"about"`
		Location          string `json:"location"`
	}

	DoctorUpdateStatusResponse struct {
		ID     string `json:"id"`
		Status bool   `json:"status"`
	}

	DoctorResponse struct {
		ID       string `json:"id"`
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
	}
)
