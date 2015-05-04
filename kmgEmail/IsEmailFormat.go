package kmgEmail

import (
	"regexp"
)

var emailFormatPregex = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)

func IsEmailFormat(email string) bool {
	return emailFormatPregex.MatchString(email)
}
