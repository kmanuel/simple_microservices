package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/transport"
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

const taskType = "screenshot"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	initMinio()

	var imageService service.ImageService
	imageService = service.NewScreenshotService()

	go startPrometheus()
	startFaktory(imageService)
}

func initMinio() {
	minioconnector.Init(
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

func startFaktory(s service.ImageService) {
	fs := service.NewFaktoryService(taskType)
	statusClient := service.NewStatusClient()

	s = middleware.NewPrometheusMiddleware(s, taskType)
	s = middleware.NewRequestStatusMiddleware(statusClient, s)
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(s))
}
