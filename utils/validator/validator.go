package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"talkspace-api/modules/user/entity"
	"talkspace-api/utils/constant"
	"time"
)

func IsDataEmpty(fields []string, data ...interface{}) error {
	if len(fields) != len(data) {
		return errors.New("column names and data length mismatch")
	}

	for i, value := range data {
		switch v := value.(type) {
		case string:
			if v == "" {
				return fmt.Errorf("%s is empty", fields[i])
			}
		case int:
			if v == 0 {
				return fmt.Errorf("%s is empty", fields[i])
			}
		case time.Time:
			if v.IsZero() {
				return fmt.Errorf("%s is empty", fields[i])
			}
		case []interface{}:
			if len(v) == 0 {
				return fmt.Errorf("%s is empty", fields[i])
			}
		case []string:
			if len(v) == 0 {
				return fmt.Errorf("%s is empty", fields[i])
			}
		case []int:
			if len(v) == 0 {
				return fmt.Errorf("%s is empty", fields[i])
			}
		default:
			if reflect.TypeOf(v).Kind() == reflect.Slice {
				slice := reflect.ValueOf(v)
				if slice.Len() == 0 {
					return fmt.Errorf("%s is empty", fields[i])
				}
			} else {
				return fmt.Errorf("unsupported data type for %s: %T", fields[i], v)
			}
		}
	}
	return nil
}

func IsEmailValid(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email format. example: emailname@gmail.com")
	}
	return nil
}

func IsMinLengthValid(minLength int, fields map[string]string) error {
	for fieldName, fieldValue := range fields {
		if len(fieldValue) < minLength {
			return fmt.Errorf("minimum length for field %s is %d characters", fieldName, minLength)
		}
	}
	return nil
}

func IsMaxLengthValid(maxLength int, fields map[string]string) error {
	for fieldName, fieldValue := range fields {
		if len(fieldValue) > maxLength {
			return fmt.Errorf("maximum length for field %s is %d characters", fieldName, maxLength)
		}
	}
	return nil
}

func IsDataValid(data interface{}, validData []interface{}, caseSensitive bool) error {
	dataStr := fmt.Sprintf("%v", data)
	validDataStr := make([]string, len(validData))
	for i, v := range validData {
		validDataStr[i] = fmt.Sprintf("%v", v)
	}

	if !caseSensitive {
		dataStr = strings.ToLower(dataStr)
		for i, v := range validDataStr {
			validDataStr[i] = strings.ToLower(v)
		}
	}

	for _, validValue := range validDataStr {
		if dataStr == validValue {
			return nil
		}
	}

	return errors.New(constant.ERROR_DATA_INVALID + strings.Join(validDataStr, ", "))
}

func IsDateValid(date string) error {
	if date == "" {
		return nil
	}

	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(date) {
		return errors.New(constant.ERROR_DATE_FORMAT)
	}

	return nil
}

func ConvertToTime(val interface{}) *time.Time {
	if val == nil {
		return nil
	}
	t := time.Unix(int64(val.(float64)), 0)
	return &t
}

func JSONReader(v interface{}) io.Reader {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(v)
	return b
}

func RemoveNilValues(slice []interface{}) []interface{} {
	result := []interface{}{}
	for _, v := range slice {
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

func GetStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func GetIntFromMap(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}

func GetInt64FromMap(m map[string]interface{}, key string) int64 {
	if v, ok := m[key].(float64); ok {
		return int64(v)
	}
	return 0
}

func MapToUserEntity(source map[string]interface{}) (entity.User, error) {
	var user entity.User

	if id, ok := source["id"].(string); ok {
		user.ID = id
	} else {
		return user, fmt.Errorf("invalid or missing field: id")
	}

	if fullname, ok := source["fullname"].(string); ok {
		user.Fullname = fullname
	}

	if email, ok := source["email"].(string); ok {
		user.Email = email
	}

	if password, ok := source["password"].(string); ok {
		user.Password = password
	}

	if newPassword, ok := source["newPassword"].(string); ok {
		user.NewPassword = newPassword
	}

	if confirmPassword, ok := source["confirmPassword"].(string); ok {
		user.ConfirmPassword = confirmPassword
	}

	if profilePicture, ok := source["profilePicture"].(string); ok {
		user.ProfilePicture = profilePicture
	}

	if birthdate, ok := source["birthdate"].(string); ok {
		user.Birthdate = birthdate
	}

	if gender, ok := source["gender"].(string); ok {
		user.Gender = gender
	}

	if bloodType, ok := source["bloodType"].(string); ok {
		user.BloodType = bloodType
	}

	if height, ok := source["height"].(float64); ok {
		user.Height = int(height)
	}

	if weight, ok := source["weight"].(float64); ok {
		user.Weight = int(weight)
	}

	if role, ok := source["role"].(string); ok {
		user.Role = role
	}

	if otp, ok := source["otp"].(string); ok {
		user.OTP = otp
	}

	if otpExpiration, ok := source["otpExpiration"].(float64); ok {
		user.OTPExpiration = int64(otpExpiration)
	}

	if createdAt, ok := source["createdAt"].(float64); ok {
		user.CreatedAt = time.Unix(int64(createdAt), 0)
	}

	if updatedAt, ok := source["updatedAt"].(float64); ok {
		user.UpdatedAt = time.Unix(int64(updatedAt), 0)
	}

	if deletedAt, ok := source["deletedAt"].(interface{}); ok {
		user.DeletedAt = ConvertToTime(deletedAt)
	}

	return user, nil
}
