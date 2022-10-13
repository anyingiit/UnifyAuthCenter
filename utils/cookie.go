package utils

import "strings"

func ParseCookieValue(cookieStr string, key string) string {
	// `target=456; a=123` -> `456; a=123`
	value := strings.Split(cookieStr, key+"=")[1]

	// `456; a=123` -> `456`
	value = strings.Split(value, ";")[0]

	return value
}
