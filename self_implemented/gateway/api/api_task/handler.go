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

type TaskHandler struct {
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

func (h *TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	h.RequestCounter.With(prometheus.Labels{"controller": "gateway", "type": "get_tasks"}).Inc()
	log.Info("received request for all tasks")

	requestServiceUrl, e := url.Parse("http://request_service:8080")
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
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

	//err := h.publishScreenshotTask("screenshot", buf.String(), task)
	err := h.publishScreenshotTask("screenshot", task)
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

func (h *TaskHandler) ServeMostSignificantHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.createMostSignificantTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) createMostSignificantTask(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("tasks.significant.create")

	task := new(MostSignificantImageTask)
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

	err := h.publishTask("most_significant_image", buf.String(), task.ID)
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

func (h *TaskHandler) ServeOptimizationHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.createOptimizationTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) createOptimizationTask(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("tasks.optimization.create")

	task := new(OptimizationTask)
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

	err := h.publishTask("optimization", buf.String(), task.ID)
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

func (h *TaskHandler) ServePortraitHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.createPortraitTask
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) createPortraitTask(w http.ResponseWriter, r *http.Request) {
	jsonapiRuntime := jsonapi.NewRuntime().Instrument("tasks.portrait.create")

	task := new(PortraitTask)
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

	err := h.publishTask("portrait", buf.String(), task.ID)
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

func (h *TaskHandler) publishTask(queue string, taskJson string, id string) error {
	h.RequestCounter.With(prometheus.Labels{"controller": "gateway", "type": "create_task"}).Inc()

	_, err := http.Post("http://request_service:8080/tasks", "application/json", bytes.NewBuffer([]byte(taskJson)))
	if err != nil {
		return err
	}

	err = publishToFaktory(queue, taskJson)
	if err != nil {
		return err
	}

	return nil
}

func (h *TaskHandler) publishScreenshotTask(queue string, task *ScreenShotTask) error {
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, task); err != nil {
		return err
	}
	taskJson := buf.String()
	err := publishToFaktory(queue, taskJson)
	if err != nil {
		return err
	}

	return updateStatus(queue, task)
}

func updateStatus(queue string, task *ScreenShotTask) error {
	taskStatus := &TaskStatus{
		TaskID: task.ID,
		TaskType: queue,
	}
	buf := new(bytes.Buffer)
	if err := jsonapi.MarshalPayload(buf, taskStatus); err != nil {
		return err
	}
	_, err := http.Post("http://request_service:8080/tasks", jsonapi.MediaType, buf)
	return err
}
