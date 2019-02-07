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
	FaktoryClient *service.FaktoryService
}

func (h *ProxyHandler) GetFaktoryInfo(w http.ResponseWriter, r *http.Request) {
	log.Info("getting faktory info")
	faktoryService := service.NewFaktoryService()

	info, err := faktoryService.Info()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	faktoryInfo := toFaktoryInfo(info)


	log.Info("got info=", faktoryInfo)

	err = jsonapi.MarshalPayload(w, faktoryInfo)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func toFaktoryInfo(info map[string]interface{}) *FaktoryInfo {
	faktoryPart := info["faktory"].(map[string]interface{})

	queues := make(map[string]float64)

	for k, v := range faktoryPart["queues"].(map[string]interface{}) {
		queues[k] = v.(float64)
	}


	return &FaktoryInfo{
		TotalProcessed: faktoryPart["total_processed"].(float64),
		TotalQueues: faktoryPart["total_queues"].(float64),
		TotalEnqueued: faktoryPart["total_enqueued"].(float64),
		TotalFailures: faktoryPart["total_failures"].(float64),
		Queues: queues,
	}
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
