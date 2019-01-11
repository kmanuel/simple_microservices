package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/service"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"controller", "status"},
	)
)

type PrometheusPublishTaskMiddleware struct {
	Next    service.FaktoryPublishService
}

func (mw PrometheusPublishTaskMiddleware) PublishTask(task *model.OptimizationTask) error  {
	requests.With(prometheus.Labels{"controller": "screenshot", "status": "fetched"}).Inc()
	err := mw.Next.PublishTask(task)
	if err != nil {
		requests.With(prometheus.Labels{"controller": "screenshot", "status": "failed"}).Inc()
	}
	return err
}

type PrometheusProcessTaskMiddleware struct {
	Next	service.OptimizationService
}

func (mw PrometheusProcessTaskMiddleware) HandleTask(task *model.OptimizationTask) error {
	requests.With(prometheus.Labels{"controller": "screenshot", "status": "processing"}).Inc()
	err := mw.Next.HandleTask(task)
	if err != nil {
		requests.With(prometheus.Labels{"controller": "screenshot", "status": "failed"}).Inc()
	}
	return err
}
