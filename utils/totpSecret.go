package utils

import "os"

// TOTPSecret returns the shared TOTP secret from UNIFYAUTH_TOTP_SECRET.
//
// It replaces a string literal that sat in handles/authCenter/handle.go in a
// public repository. That is the smaller half of the problem, and this file
// exists partly to write the larger half down where the next person will find
// it:
//
// **This service validates every user's code against one shared secret.** A
// per-user secret is enrolled once, per person, and stored with that person's
// account; a single constant shared by everyone is not a second factor at all.
// Whoever holds it can produce a valid code for any account, so the check
// proves only that the caller knows a value the server also knows — which is
// what a password already does.
//
// Moving the value into the environment stops it being published. It does not
// make the design sound. Before this service is deployed again, ValidateTOTP's
// second argument has to come from the account being authenticated, not from
// here, and this function should disappear with the change.
//
// An unset variable returns the empty string, which ValidateTOTP rejects, so a
// missing secret fails every login rather than passing it.
func TOTPSecret() string {
	return os.Getenv("UNIFYAUTH_TOTP_SECRET")
}
