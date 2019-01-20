package api_task

import (
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
	proxyTo("http://request_service:8080", w, r)
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
	log.Info("proxying to crop service")
	proxyTo("http://crop:8080", w, r)
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
	log.Info("proxying to screenshot service")
	proxyTo("http://screenshot:8080", w, r)
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
	log.Info("proxying to most_significant_image service")
	proxyTo("http://most_significant_image:8080", w, r)
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
	log.Info("proxying to optimization service")
	proxyTo("http://optimization:8080", w, r)
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
	log.Info("proxying to portrait service")
	proxyTo("http://portrait:8080", w, r)
}

func proxyTo(proxyTargetHost string, w http.ResponseWriter, r *http.Request) {
	requestServiceUrl, e := url.Parse(proxyTargetHost)
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}
