package api_task

import (
	"bytes"
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

func (h *TaskHandler) ServeScreenshotHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.createScreenshotTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) createScreenshotTask(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("tasks.screenshot.create")

	task := new(ScreenShotTask)
	task.ID = uuid.New().String()

	// unmarshal request body
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send to request service
	err := sendToRequestService(task.ID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// publish to faktory
	buf := new(bytes.Buffer)
	if err := jsonapiRuntime.MarshalPayload(buf, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = publishJsonTask("screenshot", buf.String())
	if err != nil {
		log.Error("failed to publish task to faktory", task)
		w.WriteHeader(500)
		return
	}

	// send response to caller
	w.WriteHeader(201)
	if err := jsonapiRuntime.MarshalPayload(w, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TaskHandler) ServeCropHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.createCropTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) createCropTask(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("tasks.screenshot.create")

	task := new(CropTask)
	task.ID = uuid.New().String()

	// unmarshal request body
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send to request service
	err := sendToRequestService(task.ID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// publish to faktory
	buf := new(bytes.Buffer)
	if err := jsonapiRuntime.MarshalPayload(buf, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = publishJsonTask("crop", buf.String())
	if err != nil {
		log.Error("failed to publish task to faktory", task)
		w.WriteHeader(500)
		return
	}

	// send response to caller
	w.WriteHeader(201)
	if err := jsonapiRuntime.MarshalPayload(w, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	// send to request service
	err := sendToRequestService(task.ID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// publish to faktory
	err = publishToFactory(task)
	if err != nil {
		log.Error("failed to publish task to faktory", task)
		w.WriteHeader(500)
		return
	}

	// send response to caller
	w.WriteHeader(201)
	if err := jsonapiRuntime.MarshalPayload(w, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
