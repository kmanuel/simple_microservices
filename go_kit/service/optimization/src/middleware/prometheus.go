package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/src/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/src/service"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_handle_count",
			Help: "Number of handled requests.",
		},
		[]string{"type"},
	)
)

type prometheusMiddleware struct {
	next     service.TaskService
	taskType string
}

func NewPrometheusMiddleware(next service.TaskService, taskType string) service.TaskService {
	return prometheusMiddleware{next: next, taskType: taskType}
}

func (mw prometheusMiddleware) Handle(task *model.Task) error {
	requests.With(prometheus.Labels{"type": mw.taskType}).Inc()
	return mw.next.Handle(task)
}
