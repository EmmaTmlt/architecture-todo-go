package user

import (
	"context"

	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username domain.Username) (domain.User, error)
}

type IdentifyService struct {
	UserRepository UserRepository
}

type IdentifyUserInput struct {
	Username string
	Password string
}

func (s *IdentifyService) IdentifyUser(ctx context.Context, input IdentifyUserInput) (bool, error) {
	username, err := domain.ParseUsername(input.Username)
	if err != nil {
		return false, err
	}

	password, err := domain.ParsePassword(input.Password)
	if err != nil {
		return false, err
	}

	user, err := s.UserRepository.FindByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
