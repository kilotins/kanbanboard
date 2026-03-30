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

// ProjectTag checks the project tag:
//   - 2-4 characters
//   - Uppercase letters only (A-Z)
//
// Returns an error message if invalid, or empty string if valid.
func ProjectTag(tag string) string {
	if len(tag) < 2 || len(tag) > 4 {
		return "Tag must be 2-4 characters"
	}
	for _, r := range tag {
		if r < 'A' || r > 'Z' {
			return "Tag must contain only uppercase letters (A-Z)"
		}
	}
	return ""
}
