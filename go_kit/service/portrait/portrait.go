package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/portrait/transport"
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

const taskType = "portrait"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	initMinio()

	var statusClient service.StatusClient
	statusClient = service.NewStatusClient()

	var optimizationService service.ImageService
	optimizationService = service.NewOptimizationService()

	var faktoryService service.FaktoryService
	faktoryService = service.NewFaktoryService(taskType)

	var faktoryListenService service.FaktoryListenService
	faktoryListenService = faktoryService

	go startPrometheus()
	go startFaktory(faktoryListenService, optimizationService, statusClient)
	startExternalApi(faktoryService)
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

func startFaktory(fs service.FaktoryListenService, optimizationService service.ImageService, statusClient service.StatusClient) {
	optimizationService = middleware.NewPrometheusMiddleware(optimizationService, taskType)
	optimizationService = middleware.NewRequestStatusMiddleware(statusClient, optimizationService)
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(optimizationService))
}
func startExternalApi(fs service.FaktoryPublishService) {
	requestHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs),
		transport.DecodeScreenshotTask,
		transport.EncodeResponse,
	)
	http.Handle("/", requestHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
