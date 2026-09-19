package utils

import "testing"

func TestParseCookieValue(t *testing.T) {
	got := ParseCookieValue("uuid=abc-123; admin_password=xyz", "uuid")
	if got != "abc-123" {
		t.Fatalf("ParseCookieValue() = %q, want %q", got, "abc-123")
	}

	if got := ParseCookieValue("admin_password=xyz", "uuid"); got != "" {
		t.Fatalf("ParseCookieValue() on missing key = %q, want empty string", got)
	}
}
