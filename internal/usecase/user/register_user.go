package user

import (
	"context"

	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	Save(ctx context.Context, user domain.User) error
}

type RegisterService struct {
	UserSaver UserSaver
}

type RegisterUserInput struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Avatar   *string `json:"avatar"`
}

func (s *RegisterService) RegisterUser(ctx context.Context, input RegisterUserInput) (uuid.UUID, error) {
	zero := uuid.UUID{}
	parsePass, err := domain.ParsePassword(input.Password)
	if err != nil {
		return zero, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(parsePass), bcrypt.DefaultCost)
	if err != nil {
		return zero, err
	}

	username, err := domain.ParseUsername(input.Username)
	if err != nil {
		return zero, err
	}
	userID := uuid.New()

	user := domain.User{
		ID:             userID,
		Username:       username,
		HashedPassword: domain.HashedPassword(hash),
		AvatarKey:      input.Avatar,
	}

	if err := s.UserSaver.Save(ctx, user); err != nil {
		return zero, err
	}

	return user.ID, nil
}
