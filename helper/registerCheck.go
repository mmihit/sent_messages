// helper.go
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

// CheckCliniqueRegisterData validates clinic registration data
func CheckCliniqueRegisterData(body io.ReadCloser) (CliniqueRegisterApi, error) {
	var req CliniqueRegisterApi

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("failed to read request body: ", err)
		return req, errors.New("failed to read request body")
	}

	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return req, errors.New("invalid body data")
	}

	validate := validator.New()

	// Register bcrypt hash validation
	validate.RegisterValidation("bcrypt_hash", func(fl validator.FieldLevel) bool {
		hash := fl.Field().String()
		if len(hash) != 60 {
			return false
		}
		bcryptRegex := regexp.MustCompile(`^\$2[axyb]\$\d{2}\$[./A-Za-z0-9]{53}$`)
		return bcryptRegex.MatchString(hash)
	})

	// Register owner name validation
	validate.RegisterValidation("ownername", func(fl validator.FieldLevel) bool {
		owner := fl.Field().String()
		re := regexp.MustCompile(`^[\p{L} ]+$`)
		return re.MatchString(owner)
	})

	// Validate the struct
	if err := validate.Struct(req); err != nil {
		fmt.Println("Validation error:", err)
		return req, errors.New("validation failed: check your data format")
	}

	return req, nil
}

// CheckPatientRegisterData validates patient registration data
func CheckPatientRegisterData(body io.ReadCloser) (PatientRegisterApi, error) {
	var req PatientRegisterApi

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("failed to read request body: ", err)
		return req, errors.New("failed to read request body")
	}

	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return req, errors.New("invalid body data")
	}

	validate := validator.New()

	// Register future date validation
	validate.RegisterValidation("future_date", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		return ValidateRemovalDate(dateStr, req.InsertionStentJJ, "2006-01-02")
	})

	// Validate the struct
	if err := validate.Struct(req); err != nil {
		fmt.Println("Validation error:", err)
		return req, errors.New("validation failed: check your data format")
	}

	return req, nil
}

// ValidateSurgeryDate checks if date is valid and in the future
func ValidateRemovalDate(dateStr, insertionDateStr, format string) bool {
	if dateStr == "" {
		return false
	}

	// Parse the date with the specified format
	parsedDate, err := time.Parse(format, dateStr)
	if err != nil {
		fmt.Printf("Date parse error: %v (format: %s, date: %s)\n", err, format, dateStr)
		return false
	}

	// Get today's date (start of day)
	insertionDate, err := time.Parse(format, insertionDateStr)
	if err != nil {
		fmt.Printf("Date parse error: %v (format: %s, date: %s)\n", err, format, dateStr)
		return false
	}

	// Check if date is greater than today
	if !parsedDate.After(insertionDate.Truncate(24*time.Hour).AddDate(0, 0, 3)) {
		fmt.Printf("Removal Date must be three days after insertion stent JJ: %s is not after %s\n", parsedDate.Format("2006-01-02"), insertionDate.Format("2006-01-02"))
		return false
	}

	return true
}
