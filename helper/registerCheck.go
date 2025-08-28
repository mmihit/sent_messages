package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type RegisterRequestApi struct {
	Role string `json:"role"`
}

func CheckRegisterData(body io.ReadCloser) (interface{}, string, error) {
	var role RegisterRequestApi

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, "", errors.New("failed to read request body")
	}

	err = json.Unmarshal(bodyBytes, &role)
	if err != nil {
		fmt.Println(err)
		return nil, "", errors.New("invalid body json")
	}

	validate := validator.New()

	validate.RegisterValidation("bcrypt_hash", func(fl validator.FieldLevel) bool {
		hash := fl.Field().String()
		if len(hash) != 60 {
			return false
		}
		bcryptRegex := regexp.MustCompile(`^\$2[axyb]\$\d{2}\$[./A-Za-z0-9]{53}$`)
		return bcryptRegex.MatchString(hash)
	})

	validate.RegisterValidation("ownername", func(fl validator.FieldLevel) bool {
		owner := fl.Field().String()
		re := regexp.MustCompile(`^[\p{L} ]+$`)
		return re.MatchString(owner)
	})

	var res interface{}

	switch role.Role {
	case "clinique":
		var req CliniqueRegisterApi
		err = json.Unmarshal(bodyBytes, &req)
		if err != nil {
			fmt.Println(err)
			return nil, "", errors.New("invalid body data")
		}

		res = req
		if validate.Struct(req) != nil {
			fmt.Println(validate.Struct(req))
			err = errors.New("something wrong with your data")
		}

	case "patient":
		var req PatientRegisterApi
		err = json.Unmarshal(bodyBytes, &req)
		if err != nil {
			println(err)
			return nil, "", errors.New("invalid body data")
		}

		res = req
		if validate.Struct(req) != nil {
			println(validate.Struct(req))
			err = errors.New("something wrong with your data")
		}

	default:
		err = errors.New("invalid role")
	}

	return res, role.Role, err
}

func ValidateSurgeryDate(dateStr, format string) bool {
	if dateStr == "" {
		return false
	}

	// Parse the date with the specified format
	parsedDate, err := time.Parse(format, dateStr)
	if err != nil {
		return false
	}

	// Get today's date (start of day)
	today := time.Now().Truncate(24 * time.Hour)

	// Check if date is greater than today
	if !parsedDate.After(today) {
		return false
	}

	return true
}
