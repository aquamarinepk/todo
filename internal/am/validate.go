package am

import (
	"fmt"
	"strings"
)

type (
	Validate func(entity any) error
)

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
