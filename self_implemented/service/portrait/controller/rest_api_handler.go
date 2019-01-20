package handler

import (
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/self_implemented/service/portrait/model"
	"github.com/kmanuel/simple_microservices/self_implemented/service/portrait/service"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type TaskHandler interface {
	PerformTask(w http.ResponseWriter, r *http.Request)
}

type taskHandlerImpl struct {
	statusService service.TaskStatusService
	faktoryPublishService service.FaktoryPublishService
	taskType      string
}

func NewTaskHandler(faktoryPublishService service.FaktoryPublishService, statusService service.TaskStatusService, taskType string) TaskHandler {
	return taskHandlerImpl{
		faktoryPublishService: faktoryPublishService,
		statusService: statusService,
		taskType: taskType,
	}
}

func (h taskHandlerImpl) PerformTask(w http.ResponseWriter, r *http.Request) {
	task := new(model.Task)
	task.ID = uuid.New().String()
	err := jsonapi.UnmarshalPayload(r.Body, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.statusService.NotifyAboutNewTask(task.ID, h.taskType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.faktoryPublishService.Publish(task)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := jsonapi.MarshalPayload(w, task); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
