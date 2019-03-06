package main

import (
	"flag"
	"fmt"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/src/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/src/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/src/transport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
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

const taskType = "optimization"

func main() {
	minioService := initMinio()

	var imageService service.TaskService
	imageService = service.NewTaskService(minioService)

	go startPrometheus()
	startFaktory(imageService)
}

func initMinio() *minioconnector.MinioService {
	return minioconnector.NewMinioService(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startPrometheus() {
	prometheus.MustRegister(requests)
	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}

func startFaktory(s service.TaskService) {
	fs := service.NewFaktoryService(taskType)
	s = middleware.NewPrometheusMiddleware(s, taskType)
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(s))
}
