package service

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/request_service/model"
)

type RequestStatusService interface {
	GetTaskStatusList() (model.TaskStatusList, error)
}

type RequestStatusServiceImpl struct {}

func (RequestStatusServiceImpl) GetTaskStatusList() (model.TaskStatusList, error) {
	return taskStatusListFixture(), nil
}


func taskStatusListFixture() model.TaskStatusList {
	return model.TaskStatusList{
		ID: "testid",
		Tasks: []*model.TaskStatus{
			{
				TaskID: "task1",
				Status: "new",
				TaskType: "crop",
			},
			{
				TaskID: "task2",
				Status: "new",
				TaskType: "crop",
			},
			{
				TaskID: "task3",
				Status: "completed",
				TaskType: "crop",
			},
			{
				TaskID: "task4",
				Status: "new",
				TaskType: "screenshot",
			},
		},
	}
}
