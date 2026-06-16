package validator

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/ganesh/ainyx/internal/models"
)

// dobLayout is the expected date format for date of birth.
const dobLayout = "2006-01-02"

// Validator wraps go-playground/validator with custom rules for user input.
type Validator struct {
	validate *validator.Validate
}

// New creates a Validator instance with all custom validations registered.
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateCreateUser validates a CreateUserRequest and parses the dob string.
// Returns the parsed time.Time on success, or a human-readable error.
func (v *Validator) ValidateCreateUser(req *models.CreateUserRequest) (time.Time, error) {
	if err := v.validate.Struct(req); err != nil {
		return time.Time{}, formatValidationError(err)
	}
	return v.parseDob(req.Dob)
}

// ValidateUpdateUser validates an UpdateUserRequest and parses the dob string.
func (v *Validator) ValidateUpdateUser(req *models.UpdateUserRequest) (time.Time, error) {
	if err := v.validate.Struct(req); err != nil {
		return time.Time{}, formatValidationError(err)
	}
	return v.parseDob(req.Dob)
}

// parseDob parses a date string and rejects future dates.
func (v *Validator) parseDob(dob string) (time.Time, error) {
	parsed, err := time.Parse(dobLayout, dob)
	if err != nil {
		return time.Time{}, fmt.Errorf("dob must be a valid date in YYYY-MM-DD format")
	}

	if parsed.After(time.Now()) {
		return time.Time{}, fmt.Errorf("dob cannot be in the future")
	}

	return parsed, nil
}

// formatValidationError converts validator.ValidationErrors into a single
// human-readable message for the first failing field.
func formatValidationError(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	for _, fe := range validationErrors {
		switch fe.Field() {
		case "Name":
			if fe.Tag() == "required" {
				return fmt.Errorf("name is required")
			}
			if fe.Tag() == "min" {
				return fmt.Errorf("name must be at least %s characters", fe.Param())
			}
		case "Dob":
			if fe.Tag() == "required" {
				return fmt.Errorf("dob is required")
			}
		}
	}

	return fmt.Errorf("validation failed: %s", validationErrors.Error())
}
