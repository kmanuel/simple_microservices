package api

import (
	"github.com/google/jsonapi"
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
		panic(err)
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
	var newTask NewTaskDTO
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, newTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := OpenDb()
	defer db.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskStatus := TaskStatus{
		TaskId: newTask.TaskId,
		Status: "new",
	}
	db.Create(&taskStatus)

	w.WriteHeader(http.StatusCreated)
	if err := jsonapiRuntime.MarshalPayload(w, taskStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//func (h *TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
//	h.RequestCounter.With(prometheus.Labels{"controller": "gateway", "type": "get_tasks"}).Inc()
//	log.Info("received request for all tasks")
//
//	requestServiceUrl, e := url.Parse("http://request_service:8080")
//	if e != nil {
//		panic(e)
//	}
//	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
//}
