package service

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/self_implemented/service/crop/model"
)

type FaktoryListenerService interface {
	Start() error
}

type faktoryListenerService struct {
	taskService TaskService
	taskStatusService TaskStatusService
	taskType string
}

func NewFaktoryListenerService(taskStatusService TaskStatusService, taskService TaskService, taskType string) FaktoryListenerService {
	return faktoryListenerService{
		taskService: taskService,
		taskStatusService: taskStatusService,
		taskType: taskType,
	}
}

func (s faktoryListenerService) Start() error {
	mgr := worker.NewManager()
	mgr.Concurrency = 1
	mgr.Register(s.taskType, s.handleTask)
	mgr.Queues = []string{s.taskType}
	var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
	})

	// Start processing jobs, this method does not return
	mgr.Run()

	return nil
}

func (s faktoryListenerService) handleTask(ctx worker.Context, args ...interface{}) error {
	task := new(model.Task)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		return err
	}

	err = s.taskStatusService.NotifyAboutNewTask(task.ID, s.taskType)
	if err != nil {
		return err
	}

	err = s.taskService.Handle(task)
	if err != nil {
		return err
	}

	err = s.taskStatusService.NotifyAboutCompletion(task.ID)
	if err != nil {
		return err
	}

	return nil
}
