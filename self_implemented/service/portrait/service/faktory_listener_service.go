package service

import (
	"bytes"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/kmanuel/simple_microservices/self_implemented/service/portrait/model"
)

type FaktoryListenerService interface {
	Start() error
}

type faktoryListenerService struct {
	taskService TaskService
	taskType string
}

func NewFaktoryListenerService(taskService TaskService, taskType string) FaktoryListenerService {
	return faktoryListenerService{
		taskService: taskService,
		taskType: taskType,
	}
}

func (s faktoryListenerService) Start() error {
	mgr := worker.NewManager()
	mgr.Concurrency = 1
	mgr.Register(s.taskType, s.handleTask)
	mgr.ProcessStrictPriorityQueues(s.taskType)
	var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
	})

	mgr.Run()

	return nil
}

func (s faktoryListenerService) handleTask(ctx worker.Context, args ...interface{}) error {
	task := new(model.Task)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		return err
	}

	return s.taskService.Handle(task)
}
