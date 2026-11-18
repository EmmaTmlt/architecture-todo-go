package domain

import (
	"errors"
	"unicode/utf8"
	"uuid"
)

var ErrInvalidUsername = errors.New("invalid username")
var ErrInvalidPassword = errors.New("invalid password")
var ErrUsernameTaken = errors.New("username already taken")
var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID             uuid.UUID
	Username       Username
	HashedPassword HashedPassword
	AvatarKey      *string
}

type Username string

func ParseUsername(username string) (Username, error) {
	if utf8.RuneCountInString(username) < 3 {
		return "", ErrInvalidUsername
	}
	return Username(username), nil
}

type HashedPassword string
type Password string

func ParsePassword(password string) (Password, error) {
	if utf8.RuneCountInString(password) < 3 {
		return "", ErrInvalidPassword
	}
	return Password(password), nil
}
