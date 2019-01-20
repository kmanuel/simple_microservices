package api_image

import (
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
	DispatchCounter *prometheus.CounterVec
}

func (h *ImageTaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	h.DispatchCounter.With(prometheus.Labels{"type": "get_tasks"}).Inc()
	log.Info("received request for all tasks")

	requestServiceUrl, e := url.Parse("http://request_service:8080")
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}
