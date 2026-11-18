package memory

import (
	"context"
	"sync"
	"time"

	"uuid"

	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]domain.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{tasks: make(map[uuid.UUID]domain.Task)}
}

func (r *TaskRepository) Save(ctx context.Context, task domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepository) Update(ctx context.Context, task domain.Task) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	updatedTask := domain.Task{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   &now,
		UserID:      task.UserID,
	}
	r.tasks[task.ID] = updatedTask

	return true, nil
}

func (r *TaskRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	err := ctx.Err()
	if err != nil {
		return domain.Task{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return domain.Task{}, domain.ErrTaskNotFound
}

func (r *TaskRepository) FindByName(ctx context.Context, name domain.Name) ([]domain.Task, error) {
	zero := []domain.Task{}

	err := ctx.Err()
	if err != nil {
		return zero, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := []domain.Task{}
	for _, task := range r.tasks {
		if task.Name == name {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (r *TaskRepository) FindByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error) {
	zero := []domain.Task{}

	err := ctx.Err()
	if err != nil {
		return zero, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := []domain.Task{}
	for _, task := range r.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (r *TaskRepository) FindAll(ctx context.Context) ([]domain.Task, error) {
	zero := []domain.Task{}

	err := ctx.Err()
	if err != nil {
		return zero, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := []domain.Task{}
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
