package domain

import (
	"errors"
	"time"
	"uuid"
)

var ErrInvalidName = errors.New("invalid name")
var ErrInvalidDescription = errors.New("invalid description")
var ErrInvalidStatus = errors.New("invalid status")
var ErrInvalidUpdatedAt = errors.New("invalid invalid updated at")
var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID          uuid.UUID
	Name        Name
	Description string
	Status      TaskStatus
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	UserID      uuid.UUID
}

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in progress"
	TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) Validate() error {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return nil
	default:
		return ErrInvalidStatus
	}
}

func NewTask(userID uuid.UUID, name, description string) (Task, error) {
	zero := Task{}
	parseName, err := ParseName(&name)
	if err != nil {
		return zero, err
	}
	if description == "" {
		return zero, ErrInvalidDescription
	}

	now := time.Now()
	return Task{
		ID:          uuid.New(),
		Name:        parseName,
		Description: description,
		Status:      TaskStatusTodo,
		UserID:      userID,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}, nil
}

type Name string

func ParseName(name *string) (Name, error) {
	zero := Name("")
	if name == nil {
		return "", ErrInvalidName
	}
	if Name(*name) == zero {
		return "", ErrInvalidName
	}
	return Name(*name), nil
}
