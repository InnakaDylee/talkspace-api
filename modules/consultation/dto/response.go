package dto

type RoomRes struct {
	ID   string `json:"id"`
	DoctorProfilePicture     string `json:"doctor_profile_picture"`
	UserProfilePicture	  string `json:"user_profile_picture"`
	DoctorName	  string `json:"doctor_name"`
	UserName        string `json:"user_name"`
}	

type DoctorRes struct {
	ID   string `json:"id"`
	Fullname     string `json:"fullname"`
	Email	  string `json:"email"`
	ProfilePicture        string `json:"profilePicture"`
	Role        string `json:"role"`
	Specialist        string `json:"specialist"`
	Experience        string `json:"experience"`
	Gender 	  string `json:"gender"`
	Alumnus		string `json:"alumnus"`
	AboutMe		string `json:"about_me"`
}