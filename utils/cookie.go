package utils

import "strings"

func ParseCookieValue(cookieStr string, key string) string {
	// `target=456; a=123` -> `456; a=123`
	if s1 := strings.Split(cookieStr, key+"="); len(s1) >= 2 {
		// `456` do nothing -> `456`
		// `456; a=123` -> `456`
		if s2 := strings.Split(s1[1], ";"); len(s2) >= 1 {
			return s2[0]
		}
	}
	return ""
}
