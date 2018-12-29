package api_image

import (
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_task"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type NewTaskType struct {
	Id string `json:"id"`
}

type ImageTaskHandler struct{
	RequestCounter *prometheus.CounterVec
}

func (h *ImageTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *ImageTaskHandler) HandleTaskCreation(w http.ResponseWriter, r *http.Request) {
	handler := api_task.TaskHandler{RequestCounter: h.RequestCounter}
	handler.HandleTaskCreation(w, r)
}

func (h *ImageTaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	h.RequestCounter.With(prometheus.Labels{"controller":"gateway", "type": "get_tasks"}).Inc()
	log.Info("received request for all tasks")

	requestServiceUrl, e := url.Parse("http://request_service:8080")
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}
