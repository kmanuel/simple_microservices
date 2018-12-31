package api

import (
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"net/http"
)

type TaskHandler struct {
	RequestCounter *prometheus.CounterVec
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodGet:
		methodHandler = h.getTasks
	case http.MethodPost:
		methodHandler = h.createTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	log.Info("getting all tasks")

	db, err := OpenDb()
	defer db.Close()
	if err != nil {
		log.Error("failed to open db", err)
		w.WriteHeader(500)
		return
	}

	var tasks []*TaskStatus
	if err := db.Find(&tasks).Error; err != nil {
		log.Error("failed to fetch all taskStatus from db")
		w.WriteHeader(500)
		return
	}
	list := TaskStatusList{
		ID: "1",
		Tasks: tasks,
	}

	w.WriteHeader(http.StatusCreated)
	if err := jsonapi.NewRuntime().MarshalPayload(w, &list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TaskHandler)  createTask(w http.ResponseWriter, r *http.Request) {
	log.Info("creating new TaskStatus")
	jsonapiRuntime := jsonapi.NewRuntime()

	// unmarshal request body
	newTask := new(TaskStatus)
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, newTask); err != nil {
		log.Error("unmarshalling failure ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newTask.Status = "new"

	db, err := OpenDb()
	defer db.Close()
	if err != nil {
		log.Error("OpenDb() failure ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.Create(&newTask)

	w.WriteHeader(http.StatusCreated)
	if err := jsonapiRuntime.MarshalPayload(w, newTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TaskHandler) ServeUpdateStatus(w http.ResponseWriter, r *http.Request) {
	log.Info("serving update request")
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.updateTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	log.Info("received update request")

	taskId := mux.Vars(r)["id"]

	jsonapiRuntime := jsonapi.NewRuntime()

	// unmarshal request body
	updateRequest := new(TaskStatus)
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, updateRequest); err != nil {
		log.Error("unmarshalling failure ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := OpenDb()
	defer db.Close()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var taskStatus TaskStatus
	if err := db.Where("task_id = ? ", taskId).First(&taskStatus).Error; err != nil {
		log.Error("failed to find taskStatus object to update")
		w.WriteHeader(500)
		return
	}

	taskStatus.Status = updateRequest.Status
	db.Save(&taskStatus)

	w.WriteHeader(200)
}


