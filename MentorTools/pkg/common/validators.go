package common

import "regexp"

// emailRegex is used to validate the format of email addresses.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if the provided email is in a valid format.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
