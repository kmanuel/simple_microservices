package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
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

func (mw PrometheusPublishTaskMiddleware) PublishTask(task *model.CropTask) error  {
	requests.With(prometheus.Labels{"controller": "crop", "status": "fetched"}).Inc()
	err := mw.Next.PublishTask(task)
	if err != nil {
		requests.With(prometheus.Labels{"controller": "crop", "status": "failed"}).Inc()
	}
	return err
}

type PrometheusProcessTaskMiddleware struct {
	Next	service.CropService
}

func (mw PrometheusProcessTaskMiddleware) CropImage(task *model.CropTask) (string, error) {
	requests.With(prometheus.Labels{"controller": "crop", "status": "processing"}).Inc()
	res, err := mw.Next.CropImage(task)
	if err != nil {
		requests.With(prometheus.Labels{"controller": "crop", "status": "failed"}).Inc()
	}
	return res, err
}
