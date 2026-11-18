package task

import (
	"context"
	"errors"
	"testing"

	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/adapters/memory"
	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

func TestCreateTask(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	tests := []struct {
		name    string
		input   CreateTaskInput
		wantErr error
	}{
		{
			name:    "should create a task successfully",
			input:   CreateTaskInput{UserId: userID, Name: "Task 1", Description: "descriptionnnnnnn"},
			wantErr: nil,
		},
		{
			name:    "should not create a task without a name",
			input:   CreateTaskInput{UserId: userID, Name: "", Description: "descriptionnnnnnn"},
			wantErr: domain.ErrInvalidName,
		},
		{
			name:    "should not create a task without a description",
			input:   CreateTaskInput{UserId: userID, Name: "Task 1", Description: ""},
			wantErr: domain.ErrInvalidDescription,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewTaskRepository()
			service := &CreateService{TaskSaver: repo}

			id, err := service.CreateTask(ctx, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("CreateTask() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			saved, err := repo.FindByName(ctx, domain.Name(tt.input.Name))
			if err != nil {
				t.Fatalf("task was not saved: %v", err)
			}
			if len(saved) != 1 || saved[0].ID != id {
				t.Errorf("saved tasks = %v, want exactly one with id %v", saved, id)
			}
		})
	}
}
