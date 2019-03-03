package api_task

import (
	"bytes"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/model"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/service"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	DispatchCounter *prometheus.CounterVec
	FaktoryClient   *service.FaktoryService
}

func (h *ProxyHandler) GetFaktoryInfo(w http.ResponseWriter, r *http.Request) {
	log.Info("getting faktory info")
	faktoryService := service.NewFaktoryService()

	info, err := faktoryService.Info()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Info("got info=", info)

	err = jsonapi.MarshalPayload(w, info)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProxyHandler) ProxyToRequestService(w http.ResponseWriter, r *http.Request) {
	log.Info("received request for all tasks")
	h.DispatchCounter.With(prometheus.Labels{"type": "request_service"}).Inc()
	proxyTo("http://request_service:8080", w, r)
}

func (h *ProxyHandler) CreateCropTask(w http.ResponseWriter, r *http.Request) {
	task := new(model.CropTask)
	task.ID = uuid.New().String()

	err := h.dispatchTask(task, "crop", r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProxyHandler) CreateScreenshotTask(w http.ResponseWriter, r *http.Request) {
	task := new(model.ScreenshotTask)
	task.ID = uuid.New().String()

	err := h.dispatchTask(task, "screenshot", r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProxyHandler) CreateMostSignificantImageTask(w http.ResponseWriter, r *http.Request) {
	task := new(model.MostSignificantImageTask)
	task.ID = uuid.New().String()

	err := h.dispatchTask(task, "most_significant_image", r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProxyHandler) CreateOptimizationTask(w http.ResponseWriter, r *http.Request) {
	task := new(model.OptimizationTask)
	task.ID = uuid.New().String()

	err := h.dispatchTask(task, "optimization", r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProxyHandler) CreatePortraitTask(w http.ResponseWriter, r *http.Request) {
	task := new(model.PortraitTask)
	task.ID = uuid.New().String()

	err := h.dispatchTask(task, "portrait", r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProxyHandler) dispatchTask(
	task interface{},
	taskType string,
	r *http.Request,
	w http.ResponseWriter) error {

	log.Info("dispatching task of type " + taskType)

	err := jsonapi.UnmarshalPayload(r.Body, task)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	err = jsonapi.MarshalPayload(buf, task)
	if err != nil {
		return err
	}

	faktoryService := service.NewFaktoryService()

	err = faktoryService.PublishToFaktory(taskType, buf.String())
	if err != nil {
		return err
	}

	h.DispatchCounter.With(prometheus.Labels{"type": taskType}).Inc()

	w.WriteHeader(http.StatusCreated)
	if err := jsonapi.MarshalPayload(w, task); err != nil {
		return err
	}

	return nil
}

func proxyTo(proxyTargetHost string, w http.ResponseWriter, r *http.Request) {
	requestServiceUrl, e := url.Parse(proxyTargetHost)
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}
