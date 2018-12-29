package api_task

import (
	"bytes"
	"encoding/json"
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type NewTaskType struct {
	Id string `json:"id"`
}

type TaskHandler struct{
	RequestCounter *prometheus.CounterVec
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.HandleTaskCreation
	case http.MethodGet:
		methodHandler = h.HandleGetTasks
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	h.RequestCounter.With(prometheus.Labels{"controller":"gateway", "type": "get_tasks"}).Inc()
	log.Info("received request for all tasks")

	requestServiceUrl, e := url.Parse("http://request_service:8080")
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}

func (h *TaskHandler) HandleTaskCreation(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("tasks.create")
	h.RequestCounter.With(prometheus.Labels{"controller":"gateway", "type": "create_task"}).Inc()
	log.Info("received request for new task")

	task := new(Task)
	task.ID = uuid.New().String()

	// unmarshal request body
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task.TaskParams == nil {
		task.TaskParams = make(map[string]interface{})
	}

	// send to request service
	err := sendToRequestService(task.ID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// publish to faktory
	err = publishToFactory(task)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// send response to caller
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	if err := jsonapiRuntime.MarshalPayload(w, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendToRequestService(taskId string) error {
	var nt NewTaskType
	nt.Id = taskId
	marshal, e := json.Marshal(nt)
	if e != nil {
		panic(e)
	}
	_, err := http.Post("http://request_service:8080/tasks", "application/json", bytes.NewBuffer([]byte(marshal)))
	return err
}

func publishToFactory(t *Task) error {
	log.Info("publish to faktory")
	client, err := faktory.Open()
	if err != nil {
		return err
	}
	job := faktory.NewJob(t.Type, &t.TaskParams)
	job.Queue = t.Type
	t.TaskParams["id"] = t.ID
	job.Custom = t.TaskParams
	err = client.Push(job)

	return err
}
