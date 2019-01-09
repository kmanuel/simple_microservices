package middleware

import (
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/model"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
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

func (mw PrometheusPublishTaskMiddleware) PublishTask(task *model.MostSignificantImageTask) error  {
	requests.With(prometheus.Labels{"controller": "most_significant_image", "status": "fetched"}).Inc()
	err := mw.Next.PublishTask(task)
	if err != nil {
		requests.With(prometheus.Labels{"controller": "most_significant_image", "status": "failed"}).Inc()
	}
	return err
}

type PrometheusProcessTaskMiddleware struct {
	Next	service.MostSignificantImageService
}

func (mw PrometheusProcessTaskMiddleware) ExtractMostSignificantImage(task *model.MostSignificantImageTask) (outputImagePath string, err error) {
	requests.With(prometheus.Labels{"controller": "most_significant_image", "status": "processing"}).Inc()
	res, err := mw.Next.ExtractMostSignificantImage(task)
	if err != nil {
		requests.With(prometheus.Labels{"controller": "most_significant_image", "status": "failed"}).Inc()
	}
	return res, err
}
