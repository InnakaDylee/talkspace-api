package constant

// Roles
const (
	USER   = "user"
	DOCTOR = "doctor"
	ADMIN  = "admin"
)

// Success
const (
	SUCCESS_LOGIN             = "logged in successfully"
	SUCCESS_REGISTER          = "register successfully"
	SUCCESS_CREATED           = "data created successfully"
	SUCCESS_RETRIEVED         = "data retrieved successfully"
	SUCCESS_UPDATED           = "data updated successfully"
	SUCCESS_DELETED           = "data deleted successfully"
	SUCCESS_PROFILE_UPDATED   = "profile data updated successfully"
	SUCCESS_PROFILE_RETRIEVED = "profile data retrieved successfully"
	SUCCESS_PASSWORD_UPDATED  = "password updated successfully"
	SUCCESS_OTP_SENT          = "otp sent successfully"
	SUCCESS_OTP_VERIFIED      = "otp verification successfully"
	SUCCESS_VERIFICATION      = "verification successfully"
)

// Error
const (
	ERROR_ID_NOTFOUND          = "id not found"
	ERROR_ID_INVALID           = "invalid id"
	ERROR_EMAIL_NOTFOUND       = "email not found"
	ERROR_EMAIL_FORMAT         = "invalid email format"
	ERROR_EMAIL_EXIST          = "email already exists"
	ERROR_EMAIL_UNREGISTERED   = "email not registered"
	ERROR_EMAIL_OTP            = "email or otp not found"
	ERROR_OTP_VERIFY           = "otp verification failed"
	ERROR_OTP_NOTFOUND         = "otp not found"
	ERROR_OTP_EXPIRED          = "otp has expired"
	ERROR_OTP_GENERATE         = "failed to generate otp"
	ERROR_OTP_INVALID          = "invalid otp"
	ERROR_OTP_RESET            = "failed to reset otp"
	ERROR_LOGIN                = "incorrect email or password"
	ERROR_PASSWORD_INVALID     = "invalid password"
	ERROR_PASSWORD_HASH        = "error hashing password"
	ERROR_PASSWORD_CONFIRM     = "password do not match"
	ERROR_OLDPASSWORD_INVALID  = "old password does not match"
	ERROR_DATA_NOTFOUND        = "data not found"
	ERROR_DATA_EMPTY           = "data is empty"
	ERROR_DATA_EXIST           = "data already exists"
	ERROR_DATA_TYPE            = "data type unsupported"
	ERROR_DATA_RETRIEVED       = "failed to retrieve data"
	ERROR_DATA_INVALID         = "invalid data. allowed data: "
	ERROR_FILE_EMPTY           = "file is empty"
	ERROR_DATE_FORMAT          = "invalid date format. expected format: '2000-12-30'"
	ERROR_MIN_LENGTH           = "minimum length is %d characters"
	ERROR_MAX_LENGTH           = "maximum length is %d characters"
	ERROR_TOKEN_INVALID        = "invalid token"
	ERROR_TOKEN_GENERATE       = "generate token failed"
	ERROR_TOKEN_NOTFOUND       = "token not found"
	ERROR_TOKEN_VERIFICATION   = "failed to generate token verification"
	ERROR_ACCOUNT_VERIFICATION = "failed user verification"
	ERROR_ACCOUNT_UNVERIFIED   = "account is not verified"
	ERROR_TEMPLATE_FILE        = "invalid template file"
	ERROR_TEMPLATE_READER      = "failed to read email template"
	ERROR_ROLE_ACCESS          = "not authorized to access this resource"
)
