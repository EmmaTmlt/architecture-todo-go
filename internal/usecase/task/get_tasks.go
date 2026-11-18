package task

import (
	"context"

	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

type TaskGetter interface {
	FindByName(ctx context.Context, name domain.Name) ([]domain.Task, error)
	FindAll(ctx context.Context) ([]domain.Task, error)
	FindByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error)
}

type GetTaskService struct {
	TaskGetter TaskGetter
}

type GetTaskInput struct {
	Name   *string
	Status *domain.TaskStatus
}

func (s *GetTaskService) GetTask(ctx context.Context, input GetTaskInput) ([]domain.Task, error) {
	zero := []domain.Task{}

	if input.Status != nil {
		err := input.Status.Validate()
		if err != nil {
			return zero, err
		}
		tasks, err := s.TaskGetter.FindByStatus(ctx, *input.Status)
		if err != nil {
			return zero, err
		}
		return tasks, nil
	}

	if input.Name != nil {
		name := domain.Name("")
		name, err := domain.ParseName(input.Name)
		if err != nil {
			return zero, err
		}
		tasks, err := s.TaskGetter.FindByName(ctx, name)
		if err != nil {
			return zero, err
		}
		return tasks, nil
	}

	task, err := s.TaskGetter.FindAll(ctx)
	if err != nil {
		return zero, err
	}
	return task, nil
}
