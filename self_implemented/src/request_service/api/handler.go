package api

import (
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/data"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"net/http"
)

type TaskHandler struct {
	RequestCounter *prometheus.CounterVec
	NotFoundHandler http.Handler
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	log.Info("getting all tasks")

	list, err := data.FetchTaskList()
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := jsonapi.MarshalPayload(w, &list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	log.Info("creating new model.TaskStatus")

	newTask := new(model.TaskStatus)
	if err := jsonapi.UnmarshalPayload(r.Body, newTask); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := data.CreateNewTask(newTask); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := jsonapi.MarshalPayload(w, newTask); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	log.Info("received update request")

	updateRequest := new(model.TaskStatus)
	if err := jsonapi.UnmarshalPayload(r.Body, updateRequest); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updateRequest.TaskID = mux.Vars(r)["id"]

	if err := data.UpdateTaskStatus(updateRequest); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := jsonapi.MarshalPayload(w, updateRequest); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
