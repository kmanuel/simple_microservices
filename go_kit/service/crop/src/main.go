package main

import (
	"flag"
	"fmt"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/src/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/src/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/src/transport"
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

var taskType = "crop"

func main() {
	minioService := initMinio()
	go startPrometheus()
	startFaktory(minioService)
}

func initMinio() *minioconnector.MinioService {
	return minioconnector.NewMinioService(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("INPUT_BUCKET_NAME"),
		taskType)
}

func startPrometheus() {
	prometheus.MustRegister(requests)
	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}

func startFaktory(minioService *minioconnector.MinioService) {
	fs := service.NewFaktoryService(taskType)
	cropService := service.NewCropService(minioService)

	cropService = middleware.NewPrometheusMiddleware(cropService, taskType)
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(cropService))
}
