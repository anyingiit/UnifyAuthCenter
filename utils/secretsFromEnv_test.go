package utils

import "testing"

// The property worth testing here is not "the right password passes" -- it is
// "a missing secret fails". An unset environment variable used to be the shape
// of bug that turns a credential check into a formality: compare against "" and
// every empty submission succeeds.

func TestAdminPasswordFailsClosedWhenUnset(t *testing.T) {
	t.Setenv("UNIFYAUTH_ADMIN_PASSWORD", "")
	if IsValidateAdminiPassword("") {
		t.Fatal("an unset admin password accepted an empty submission")
	}
	if IsValidateAdminiPassword("anything") {
		t.Fatal("an unset admin password accepted a submission")
	}
}

func TestAdminPasswordMatchesOnlyTheConfiguredValue(t *testing.T) {
	t.Setenv("UNIFYAUTH_ADMIN_PASSWORD", "configured-for-this-test-only")
	if !IsValidateAdminiPassword("configured-for-this-test-only") {
		t.Fatal("the configured admin password was rejected")
	}
	if IsValidateAdminiPassword("configured-for-this-test-onlX") {
		t.Fatal("a wrong admin password was accepted")
	}
	if IsValidateAdminiPassword("") {
		t.Fatal("an empty admin password was accepted")
	}
}

func TestInternalPasswordFailsClosedWhenUnset(t *testing.T) {
	t.Setenv("UNIFYAUTH_INTERNAL_PASSWORD", "")
	if IsValidateInternalPassword("") || IsValidateInternalPassword("anything") {
		t.Fatal("an unset internal password accepted a submission")
	}
}

func TestInternalPasswordMatchesOnlyTheConfiguredValue(t *testing.T) {
	t.Setenv("UNIFYAUTH_INTERNAL_PASSWORD", "internal-for-this-test-only")
	if !IsValidateInternalPassword("internal-for-this-test-only") {
		t.Fatal("the configured internal password was rejected")
	}
	if IsValidateInternalPassword("internal-for-this-test-onlX") {
		t.Fatal("a wrong internal password was accepted")
	}
}

func TestTOTPSecretIsEmptyWhenUnset(t *testing.T) {
	t.Setenv("UNIFYAUTH_TOTP_SECRET", "")
	if TOTPSecret() != "" {
		t.Fatal("TOTPSecret invented a secret")
	}
	// ValidateTOTP must reject anything against an empty secret, so that a
	// missing configuration fails every login rather than passing it.
	if ValidateTOTP("000000", TOTPSecret()) {
		t.Fatal("an empty TOTP secret validated a passcode")
	}
}
