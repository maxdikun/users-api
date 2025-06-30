package entities

import (
	"fmt"
	"slices"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type Password string

func RawPassword(value string) Password {
	return Password(value)
}

func NewPassword(value string) (Password, error) {
	if len(value) < 6 {
		return "", newValidationError("password", "should be at least 6 characters long")
	}

	charTypeCount := 0
	if slices.ContainsFunc([]rune(value), unicode.IsLower) {
		charTypeCount++
	}

	if slices.ContainsFunc([]rune(value), unicode.IsUpper) {
		charTypeCount++
	}

	if slices.ContainsFunc([]rune(value), unicode.IsDigit) {
		charTypeCount++
	}

	if charTypeCount < 2 {
		return "", newValidationError("password", "should contain at least two types of characters: lowercase, uppercase, digits")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash from password: %w", err)
	}

	return Password(hashed), nil
}

func (p Password) Compare(other string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(other)) == nil
}
