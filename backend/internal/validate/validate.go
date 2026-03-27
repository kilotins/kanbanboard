// Package validate provides input validation rules.
package validate

import "unicode"

// Password checks the password against the password policy:
//   - Minimum 8 characters
//   - At least one letter (uppercase or lowercase)
//   - At least one number
//
// Returns an error message if invalid, or empty string if valid.
func Password(password string) string {
	if len(password) < 8 {
		return "Password must be at least 8 characters"
	}

	hasLetter := false
	hasNumber := false
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasNumber = true
		}
	}

	if !hasLetter {
		return "Password must contain at least one letter"
	}
	if !hasNumber {
		return "Password must contain at least one number"
	}

	return ""
}
