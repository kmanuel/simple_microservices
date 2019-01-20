package handler

import (
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/self_implemented/service/most_significant_image/model"
	"github.com/kmanuel/simple_microservices/self_implemented/service/most_significant_image/service"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type TaskHandler interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
}

type taskHandlerImpl struct {
	statusService service.TaskStatusService
	taskService service.TaskService
	taskType      string
}

func NewTaskHandler(taskService service.TaskService, statusService service.TaskStatusService, taskType string) TaskHandler {
	return taskHandlerImpl{
		taskService: taskService,
		statusService: statusService,
		taskType: taskType,
	}
}

func (h taskHandlerImpl) ServeHttp(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.handleIncomingTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h taskHandlerImpl) handleIncomingTask(w http.ResponseWriter, r *http.Request) {
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

	err = h.taskService.Handle(task)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.statusService.NotifyAboutCompletion(task.ID)
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
