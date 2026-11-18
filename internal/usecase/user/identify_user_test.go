package user

import (
	"context"
	"errors"
	"testing"

	"github.com/EmmaTmlt/architecture-todo-go/internal/adapters/memory"
	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func TestIdentifyUser(t *testing.T) {
	ctx := context.Background()

	userRepo := memory.NewUserRepository()
	registerService := &RegisterService{UserSaver: userRepo}
	_, err := registerService.RegisterUser(ctx, RegisterUserInput{
		Username: "emmmaaa",
		Password: "zelrjkzrejanlazekd",
	})
	if err != nil {
		t.Fatalf("setup: RegisterUser() error = %v", err)
	}

	tests := []struct {
		name    string
		input   IdentifyUserInput
		want    bool
		wantErr error
	}{
		{
			name: "should identify a user with username and password successfully",
			input: IdentifyUserInput{
				Username: "emmmaaa",
				Password: "zelrjkzrejanlazekd",
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "should not identify a user when username is too short",
			input: IdentifyUserInput{
				Username: "ma",
				Password: "zelrjkzrejanlazekd",
			},
			wantErr: domain.ErrInvalidUsername,
		},
		{
			name: "should not identify a user when password is too short",
			input: IdentifyUserInput{
				Username: "emmmaaa",
				Password: "r",
			},
			wantErr: domain.ErrInvalidPassword,
		},
		{
			name: "should not identify an unknown user",
			input: IdentifyUserInput{
				Username: "unknown",
				Password: "zelrjkzrejanlazekd",
			},
			wantErr: domain.ErrUserNotFound,
		},
		{
			name: "should not identify a user with a wrong password",
			input: IdentifyUserInput{
				Username: "emmmaaa",
				Password: "wrong-password",
			},
			wantErr: bcrypt.ErrMismatchedHashAndPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identifyService := &IdentifyService{UserRepository: userRepo}

			ok, err := identifyService.IdentifyUser(ctx, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("IdentifyUser() error = %v, want %v", err, tt.wantErr)
			}
			if ok != tt.want {
				t.Errorf("IdentifyUser() ok = %v, want %v", ok, tt.want)
			}
		})
	}
}
