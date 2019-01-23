package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/transport"
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
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	initMinio()

	var statusClient service.StatusClient
	statusClient = service.NewStatusClient()

	var imageService service.ImageService
	imageService = service.NewOptimizationService()

	go startPrometheus()
	startRestApi(statusClient, imageService)
}

func initMinio() {
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startRestApi(statusClient service.StatusClient, imageService service.ImageService) {
	imageService = middleware.NewPrometheusMiddleware(imageService, taskType)
	imageService = middleware.NewRequestStatusMiddleware(statusClient, imageService)

	requestHandler := httptransport.NewServer(
		transport.CreateRestHandler(imageService),
		transport.DecodeRequest,
		transport.EncodeResponse,
	)
	http.Handle("/", requestHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func startPrometheus() {
	prometheus.MustRegister(requests)
	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}
