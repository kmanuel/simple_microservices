package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/src/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/src/service"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"type"},
	)
)

type prometheusMiddleware struct {
	next     service.ImageService
	taskType string
}

func NewPrometheusMiddleware(next service.ImageService, taskType string) service.ImageService {
	return prometheusMiddleware{next: next, taskType: taskType}
}

func (mw prometheusMiddleware) HandleTask(task *model.Task) error {
	requests.With(prometheus.Labels{"type": mw.taskType}).Inc()
	return mw.next.HandleTask(task)
}
