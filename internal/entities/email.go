package entities

import "net/mail"

type Email string

func NewEmail(value string) (Email, error) {
	addr, err := mail.ParseAddress(value)
	if err != nil {
		return "", newValidationError("email", "email is mailformed")
	}

	return Email(addr.Address), nil
}
