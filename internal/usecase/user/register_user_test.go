package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/adapters/memory"
	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		input   RegisterUserInput
		wantErr error
	}{
		{
			name: "should create a user with username and password successfully",
			input: RegisterUserInput{
				Username: "emmmaaa",
				Password: "zelrjkzrejanlazekd",
			},
			wantErr: nil,
		},
		{
			name: "should not create a user when username is too short",
			input: RegisterUserInput{
				Username: "ma",
				Password: "zelrjkzrejanlazekd",
			},
			wantErr: domain.ErrInvalidUsername,
		},
		{
			name: "should not create a user when password is too short",
			input: RegisterUserInput{
				Username: "emmmaaa",
				Password: "r",
			},
			wantErr: domain.ErrInvalidPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerService := &RegisterService{UserSaver: memory.NewUserRepository()}

			id, err := registerService.RegisterUser(ctx, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("RegisterUser() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if id != (uuid.UUID{}) {
					t.Errorf("RegisterUser() id = %v, want zero UUID on error", id)
				}
				fmt.Print("YOUHOU")
				return
			}
			if id == (uuid.UUID{}) {
				t.Error("RegisterUser() returned a zero UUID, want a generated one")
			}
			fmt.Print("YOUHOU")
		})
	}
}

func TestRegisterUser_UsernameTaken(t *testing.T) {
	ctx := context.Background()
	registerService := &RegisterService{UserSaver: memory.NewUserRepository()}
	input := RegisterUserInput{Username: "emmmaaa", Password: "zelrjkzrejanlazekd"}

	if _, err := registerService.RegisterUser(ctx, input); err != nil {
		t.Fatalf("first RegisterUser() error = %v", err)
	}

	_, err := registerService.RegisterUser(ctx, input)
	if !errors.Is(err, domain.ErrUsernameTaken) {
		t.Fatalf("second RegisterUser() error = %v, want %v", err, domain.ErrUsernameTaken)
	}
}
