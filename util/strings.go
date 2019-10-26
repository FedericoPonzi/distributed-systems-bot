package util

import (
	"regexp"
	"strings"

	"github.com/marksalpeter/token"
)

// RandString returns a random string of a given length
func RandString(maxChars int) string {
	return token.New(maxChars).Encode()
}

// FormatTitleForUrl validates and formats a string to restrictions
func FormatTitleForUrl(title string, maxlength int) (url string) {
	url = strings.ToLower(title)
	if len(url) > maxlength {
		url = url[:maxlength-1]
	}
	reg := regexp.MustCompile("[^a-zA-Z ]*")
	url = reg.ReplaceAllString(url, "")
	url = strings.Replace(url, " ", "-", -1)
	return url
}
