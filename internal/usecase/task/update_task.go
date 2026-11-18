package task

import (
	"context"
	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

type TaskRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (domain.Task, error)
	Update(ctx context.Context, task domain.Task) (bool, error)
}

type UpdateService struct {
	TaskRepository TaskRepository
}

type UpdateTaskInput struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Status      *domain.TaskStatus
}

func (s *UpdateService) UpdateTask(ctx context.Context, input UpdateTaskInput) (bool, error) {
	task, err := s.TaskRepository.FindById(ctx, input.ID)
	if err != nil {
		return false, err
	}

	status := task.Status
	if input.Status != nil {
		err = input.Status.Validate()
		if err != nil {
			return false, err
		}
		status = *input.Status
	}

	name := task.Name
	if input.Name != nil {
		name, err = domain.ParseName(input.Name)
		if err != nil {
			return false, err
		}
	}

	description := task.Description
	if input.Description != nil {
		description = *input.Description
	}

	ok, err := s.TaskRepository.Update(ctx, domain.Task{
		ID:          task.ID,
		Name:        name,
		Status:      status,
		Description: description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		UserID:      task.UserID,
	})
	return ok, err
}
