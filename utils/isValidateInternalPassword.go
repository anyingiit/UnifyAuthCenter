package utils

import (
	"crypto/subtle"
	"os"
)

// IsValidateInternalPassword reports whether validateTarget is the password
// internal services present to this one.
//
// Same change, and the same reasoning, as IsValidateAdminiPassword: the value
// was a literal in this file in a public repository, it now comes from
// UNIFYAUTH_INTERNAL_PASSWORD, an unset variable fails every check rather than
// passing it, and the comparison is constant time.
func IsValidateInternalPassword(validateTarget string) bool {
	expected := os.Getenv("UNIFYAUTH_INTERNAL_PASSWORD")
	if expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(validateTarget), []byte(expected)) == 1
}
