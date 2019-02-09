package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/transport"
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
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	initMinio()
	go startPrometheus()
	startFaktory()
}

func initMinio() {
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startFaktory() {
	fs := service.NewFaktoryService(taskType)
	cropService := service.NewCropService()
	statusClient := service.NewStatusClient()

	cropService = middleware.NewPrometheusMiddleware(cropService, taskType)
	cropService = middleware.NewRequestStatusMiddleware(statusClient, cropService)
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(cropService))
}

func startPrometheus() {
	prometheus.MustRegister(requests)
	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}
