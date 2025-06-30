package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	id               uuid.UUID
	username         Username
	email            Email
	password         Password
	createdAt        time.Time
	emailConfirmedAt *time.Time
	updatedAt        time.Time
	isDeleted        bool
}

func (u *User) Email() Email {
	return u.email
}

func LoadUser(
	id uuid.UUID,
	username Username,
	email Email,
	password Password,
	createdAt time.Time,
	emailConfirmedAt *time.Time,
	updatedAt time.Time,
	isDeleted bool,
) User {
	return User{
		id:               id,
		username:         username,
		email:            email,
		password:         password,
		createdAt:        createdAt,
		emailConfirmedAt: emailConfirmedAt,
		updatedAt:        updatedAt,
		isDeleted:        isDeleted,
	}
}

func NewUser(username Username, email Email, password Password) User {
	return User{
		id:               uuid.New(),
		username:         username,
		email:            email,
		password:         password,
		createdAt:        time.Now(),
		updatedAt:        time.Now(),
		isDeleted:        false,
		emailConfirmedAt: nil,
	}
}
