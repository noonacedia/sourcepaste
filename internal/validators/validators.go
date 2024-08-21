package validators

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// regex pattern is recommended by W3C and Web Hypertext Application Technology Working Group
var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, maxChars int) bool {
	return utf8.RuneCountInString(value) <= maxChars
}

func MinChars(value string, minChars int) bool {
	return utf8.RuneCountInString(value) >= minChars
}

func PermittedInt(permittedValues []int, value int) bool {
	return slices.Contains(permittedValues, value)
}

func EmailValid(email string, regex *regexp.Regexp) bool {
	return regex.MatchString(email)
}
