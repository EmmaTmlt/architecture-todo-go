package task

import (
	"context"

	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

type TaskSaver interface {
	Save(ctx context.Context, task domain.Task) error
}

type CreateService struct {
	TaskSaver TaskSaver
}

type CreateTaskInput struct {
	UserId      uuid.UUID
	Name        string
	Description string
}

func (s *CreateService) CreateTask(ctx context.Context, input CreateTaskInput) (uuid.UUID, error) {
	zero := uuid.UUID{}

	newTask, err := domain.NewTask(input.UserId, input.Name, input.Description)
	if err != nil {
		return zero, err
	}

	if err := s.TaskSaver.Save(ctx, newTask); err != nil {
		return zero, err
	}

	return newTask.ID, nil
}
