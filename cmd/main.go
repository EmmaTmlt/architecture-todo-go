package main

import (
	"context"
	"fmt"
	"log"

	"github.com/EmmaTmlt/architecture-todo-go/internal/adapters/memory"
	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
	"github.com/EmmaTmlt/architecture-todo-go/internal/usecase/task"
	"github.com/EmmaTmlt/architecture-todo-go/internal/usecase/user"
)

func main() {
	ctx := context.Background()

	userRepo := memory.NewUserRepository()
	taskRepo := memory.NewTaskRepository()

	registerService := &user.RegisterService{UserSaver: userRepo}
	identifyService := &user.IdentifyService{UserRepository: userRepo}
	updateService := &task.UpdateService{TaskRepository: taskRepo}
	createService := &task.CreateService{TaskSaver: taskRepo}
	getTaskService := task.GetTaskService{TaskGetter: taskRepo}

	id, err := registerService.RegisterUser(ctx, user.RegisterUserInput{
		Username: "emma",
		Password: "motdepasse-super-secret",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("utilisateur créé :", id)

	ok, err := identifyService.IdentifyUser(ctx, user.IdentifyUserInput{
		Username: "emma",
		Password: "motdepasse-super-secret",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("utilisateur sign up :", ok)

	idTask, err := createService.CreateTask(ctx, task.CreateTaskInput{
		Name:        "Task 1",
		Description: "ejanzekrjbazkerj",
		UserId:      id,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("task créé :", idTask)

	idTask2, err := createService.CreateTask(ctx, task.CreateTaskInput{
		Name:        "Task 2",
		Description: "kazejrzkjeqn",
		UserId:      id,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("task créé :", idTask2)

	name := "Task 1"
	description := "ejanzekrjbazkerj"
	status := domain.TaskStatusInProgress
	ok, err = updateService.UpdateTask(ctx, task.UpdateTaskInput{
		ID:          idTask,
		Name:        &name,
		Description: &description,
		Status:      &status,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("task mise à jour :", ok)

	tasks, err := getTaskService.GetTask(ctx, task.GetTaskInput{
		Name:   nil,
		Status: nil,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, task := range tasks {
		fmt.Println(task)
	}
}
