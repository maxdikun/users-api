package entities

import (
	"strings"
	"unicode"
)

type Username string

func NewUsername(value string) (Username, error) {
	if len(value) < 3 {
		return "", newValidationError("username", "should be at least 3 characters long")
	}

	if strings.ContainsFunc(value, unicode.IsSpace) {
		return "", newValidationError("username", "should not contain spaces")
	}

	return Username(value), nil
}
