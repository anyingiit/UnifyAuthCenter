package utils

import (
	"crypto/subtle"
	"os"
)

// IsValidateAdminiPassword reports whether validateTarget is the admin password.
//
// The password used to be a string literal compared with `==` right here, in a
// public repository: anyone who read this file knew it. It now comes from the
// UNIFYAUTH_ADMIN_PASSWORD environment variable, and the published value has
// been retired.
//
// Two deliberate properties:
//
//   - An unset or empty UNIFYAUTH_ADMIN_PASSWORD makes every check fail. A
//     missing secret must never turn into an open door, which is what returning
//     true (or comparing against "") would do.
//   - subtle.ConstantTimeCompare, not `==`. Go's string comparison returns as
//     soon as two bytes differ, so how long it takes leaks how much of the
//     password a guess got right.
func IsValidateAdminiPassword(validateTarget string) bool {
	expected := os.Getenv("UNIFYAUTH_ADMIN_PASSWORD")
	if expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(validateTarget), []byte(expected)) == 1
}
