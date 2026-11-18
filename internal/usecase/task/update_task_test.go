package task

import (
	"context"
	"errors"
	"testing"

	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/adapters/memory"
	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

func TestUpdateTask(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		input    func(id uuid.UUID) UpdateTaskInput
		wantErr  error
		wantTask domain.Task
	}{
		{
			name: "should update the status only",
			input: func(id uuid.UUID) UpdateTaskInput {
				return UpdateTaskInput{ID: id, Status: new(domain.TaskStatusDone)}
			},
			wantTask: domain.Task{Name: "Task 1", Description: "a description", Status: domain.TaskStatusDone},
		},
		{
			name: "should update the name only",
			input: func(id uuid.UUID) UpdateTaskInput {
				return UpdateTaskInput{ID: id, Name: new("Task 2")}
			},
			wantTask: domain.Task{Name: "Task 2", Description: "a description", Status: domain.TaskStatusTodo},
		},
		{
			name: "should update the description only",
			input: func(id uuid.UUID) UpdateTaskInput {
				return UpdateTaskInput{ID: id, Description: new("new description")}
			},
			wantTask: domain.Task{Name: "Task 1", Description: "new description", Status: domain.TaskStatusTodo},
		},
		{
			name: "should update every field",
			input: func(id uuid.UUID) UpdateTaskInput {
				return UpdateTaskInput{
					ID:          id,
					Name:        new("Task 2"),
					Description: new("new description"),
					Status:      new(domain.TaskStatusInProgress),
				}
			},
			wantTask: domain.Task{Name: "Task 2", Description: "new description", Status: domain.TaskStatusInProgress},
		},
		{
			name: "should not update a task that does not exist",
			input: func(uuid.UUID) UpdateTaskInput {
				return UpdateTaskInput{ID: uuid.New(), Status: new(domain.TaskStatusDone)}
			},
			wantErr: domain.ErrTaskNotFound,
		},
		{
			name: "should not update with an invalid status",
			input: func(id uuid.UUID) UpdateTaskInput {
				return UpdateTaskInput{ID: id, Status: new(domain.TaskStatus("unknown"))}
			},
			wantErr: domain.ErrInvalidStatus,
		},
		{
			name: "should not update with an empty name",
			input: func(id uuid.UUID) UpdateTaskInput {
				return UpdateTaskInput{ID: id, Name: new("")}
			},
			wantErr: domain.ErrInvalidName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewTaskRepository()
			original, err := domain.NewTask(uuid.New(), "Task 1", "a description")
			if err != nil {
				t.Fatalf("setup: NewTask() error = %v", err)
			}
			if err := repo.Save(ctx, original); err != nil {
				t.Fatalf("setup: Save() error = %v", err)
			}
			service := &UpdateService{TaskRepository: repo}

			ok, err := service.UpdateTask(ctx, tt.input(original.ID))

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("UpdateTask() error = %v, want %v", err, tt.wantErr)
			}

			got, findErr := repo.FindById(ctx, original.ID)
			if findErr != nil {
				t.Fatalf("FindById() error = %v", findErr)
			}

			if tt.wantErr != nil {
				if ok {
					t.Error("UpdateTask() ok = true, want false")
				}
				if got.Name != original.Name || got.Description != original.Description || got.Status != original.Status {
					t.Errorf("task was modified despite the error: %+v", got)
				}
				return
			}

			if !ok {
				t.Error("UpdateTask() ok = false, want true")
			}
			if got.Name != tt.wantTask.Name || got.Description != tt.wantTask.Description || got.Status != tt.wantTask.Status {
				t.Errorf("task = %+v, want name=%q description=%q status=%q",
					got, tt.wantTask.Name, tt.wantTask.Description, tt.wantTask.Status)
			}
			if got.UserID != original.UserID {
				t.Errorf("UserID changed: %v, want %v", got.UserID, original.UserID)
			}
			if got.UpdatedAt.Before(*got.CreatedAt) {
				t.Errorf("UpdatedAt %v is before CreatedAt %v", got.UpdatedAt, got.CreatedAt)
			}
		})
	}
}
