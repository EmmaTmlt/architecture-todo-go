package task

import (
	"context"
	"errors"
	"testing"

	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/adapters/memory"
	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

func TestGetTask(t *testing.T) {
	ctx := context.Background()

	repo := memory.NewTaskRepository()
	userID := uuid.New()
	seed := []struct {
		name   string
		status domain.TaskStatus
	}{
		{"Task 1", domain.TaskStatusTodo},
		{"Task 2", domain.TaskStatusTodo},
		{"Task 3", domain.TaskStatusDone},
		{"Task 3", domain.TaskStatusInProgress},
	}
	for _, s := range seed {
		task, err := domain.NewTask(userID, s.name, "a description")
		if err != nil {
			t.Fatalf("setup: NewTask() error = %v", err)
		}
		task.Status = s.status

		err = repo.Save(ctx, task)
		if err != nil {
			t.Fatalf("setup: Save() error = %v", err)
		}
	}

	tests := []struct {
		name      string
		input     GetTaskInput
		wantCount int
		wantErr   error
	}{
		{
			name:      "should return every task without filter",
			input:     GetTaskInput{},
			wantCount: 4,
		},
		{
			name:      "should filter by status",
			input:     GetTaskInput{Status: new(domain.TaskStatusTodo)},
			wantCount: 2,
		},
		{
			name:      "should filter by name",
			input:     GetTaskInput{Name: new("Task 1")},
			wantCount: 1,
		},
		{
			name:      "should return every task sharing a name",
			input:     GetTaskInput{Name: new("Task 3")},
			wantCount: 2,
		},
		{
			name:      "should return an empty list when no task matches the name",
			input:     GetTaskInput{Name: new("unknown")},
			wantCount: 0,
		},
		{
			name:    "should not filter with an invalid status",
			input:   GetTaskInput{Status: new(domain.TaskStatus("unknown"))},
			wantErr: domain.ErrInvalidStatus,
		},
		{
			name:    "should not filter with an empty name",
			input:   GetTaskInput{Name: new("")},
			wantErr: domain.ErrInvalidName,
		},
		{
			name:      "should give priority to the status over the name",
			input:     GetTaskInput{Name: new("Task 1"), Status: new(domain.TaskStatusDone)},
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &GetTaskService{TaskGetter: repo}

			tasks, err := service.GetTask(ctx, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("GetTask() error = %v, want %v", err, tt.wantErr)
			}
			if len(tasks) != tt.wantCount {
				t.Errorf("GetTask() returned %d tasks, want %d", len(tasks), tt.wantCount)
			}
		})
	}
}
