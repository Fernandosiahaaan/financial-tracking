package validations

import (
	"fmt"
	"regexp"
	"service-user/internal/model"
	"service-user/internal/model/request"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidationCreateUser(params *request.CreateUserRequest) (errRedaksi, errSystem error) {
	err := validation.ValidateStruct(params,
		validation.Field(&params.ReqID, validation.Required),
		validation.Field(&params.Username, validation.Required),
		validation.Field(&params.Password, validation.Required),
		validation.Field(&params.FullName, validation.Required),
		validation.Field(&params.Email, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("Data 'ReqId', 'Username', 'Password', 'FullName', 'Email', 'Role' must fill."), fmt.Errorf("failed validation struct request kuota cabang. err : %v", err)
	}

	if !isValidEmail(params.Email) {
		return fmt.Errorf("Data 'Email' is not valid format '%s'", params.Email), fmt.Errorf("failed email format '%s'. err : %v", params.Email, err)
	}

	if !isValidPhoneNumber(params.PhoneNumber) {
		return fmt.Errorf("Data 'phone_number' is not valid format '%s'", params.PhoneNumber), fmt.Errorf("failed phone_number format '%s'. err : %v", params.PhoneNumber, err)
	}

	if params.Role == "" {
		params.Role = model.RoleUser
	}

	return nil, nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(strings.ToLower(email))
}

func isValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\+?[0-9]{9,15}$`) // bisa disesuaikan
	return re.MatchString(phone)
}
