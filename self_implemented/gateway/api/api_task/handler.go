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
	case http.MethodGet:
		methodHandler = h.HandleGetTasks
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
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
	// unmarshal request body
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	task.ID = uuid.New().String()
	buf := new(bytes.Buffer)
	if err := jsonapiRuntime.MarshalPayload(buf, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := h.publishTask("screenshot", buf.String(), task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send response to caller
	w.WriteHeader(http.StatusCreated)
	if err := jsonapiRuntime.MarshalPayload(w, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h *TaskHandler) createCropTask(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("tasks.crop.create")

	task := new(CropTask)
	task.ID = uuid.New().String()

	// unmarshal request body
	if err := jsonapiRuntime.UnmarshalPayload(r.Body, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	if err := jsonapiRuntime.MarshalPayload(buf, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := h.publishTask("crop", buf.String(), task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send response to caller
	w.WriteHeader(http.StatusCreated)
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

func (h *TaskHandler) publishTask(queue string, taskJson string, id string) error {
	h.RequestCounter.With(prometheus.Labels{"controller":"gateway", "type": "create_task"}).Inc()
	err := sendToRequestService(id)
	if err != nil {
		return err
	}

	err = publishToFaktory(queue, taskJson)
	if err != nil {
		return err
	}

	return nil
}
