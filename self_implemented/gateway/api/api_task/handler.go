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
}

func (h *ProxyHandler) ProxyToRequestService(w http.ResponseWriter, r *http.Request) {

	log.Info("received request for all tasks")
	h.DispatchCounter.With(prometheus.Labels{"type": "request_service"}).Inc()
	proxyTo("http://request_service:8080", w, r)
}

func (h *ProxyHandler) CreateCropTask(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to crop service")

	task := new(model.CropTask)
	task.ID = uuid.New().String()
	err := jsonapi.UnmarshalPayload(r.Body, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err = jsonapi.MarshalPayload(buf, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	faktoryService := service.NewFaktoryService()

	err = faktoryService.PublishToFaktory("crop", buf.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.DispatchCounter.With(prometheus.Labels{"type": "crop"}).Inc()

	w.WriteHeader(http.StatusCreated)
	if err := jsonapi.MarshalPayload(w, task); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProxyHandler) CreateScreenshotTask(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to screenshot service")
	h.DispatchCounter.With(prometheus.Labels{"type": "screenshot"}).Inc()
	proxyTo("http://screenshot:8080", w, r)
}

func (h *ProxyHandler) CreateMostSignificantImageTask(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to most_significant_image service")
	h.DispatchCounter.With(prometheus.Labels{"type": "most_significant_image"}).Inc()
	proxyTo("http://most_significant_image:8080", w, r)
}

func (h *ProxyHandler) CreateOptimizationTask(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to optimization service")
	h.DispatchCounter.With(prometheus.Labels{"type": "optimization"}).Inc()
	proxyTo("http://optimization:8080", w, r)
}

func (h *ProxyHandler) CreatePortraitTask(w http.ResponseWriter, r *http.Request) {
	log.Info("proxying to portrait service")
	h.DispatchCounter.With(prometheus.Labels{"type": "portrait"}).Inc()
	proxyTo("http://portrait:8080", w, r)
}

func proxyTo(proxyTargetHost string, w http.ResponseWriter, r *http.Request) {
	requestServiceUrl, e := url.Parse(proxyTargetHost)
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}
