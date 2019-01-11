package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/optimization/status_client"
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

	var statusClient status_client.StatusClient
	statusClient = status_client.NewStatusClient(taskType)

	var optimizationService service.OptimizationService
	optimizationService = service.NewOptimizationService()

	var faktoryService service.FaktoryService
	faktoryService = service.NewFaktoryService(taskType)

	go startPrometheus()

	var faktoryListenService service.FaktoryListenService
	faktoryListenService = faktoryService
	go startFaktory(faktoryListenService, optimizationService, statusClient)

	startExternalApi(statusClient, faktoryService)
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
	var addr = flag.String("listen-address", ":"+os.Getenv("OPTIMIZATION_SERVICE_PROMETHEUS_PORT"), "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}

func startFaktory(fs service.FaktoryListenService, optimizationService service.OptimizationService, statusClient status_client.StatusClient) {
	optimizationService = middleware.StatusPerformMiddleware{StatusClient: statusClient, Next: optimizationService}
	optimizationService = middleware.PrometheusProcessTaskMiddleware{Next: optimizationService}
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(optimizationService))
}
func startExternalApi(statusClient status_client.StatusClient, fs service.FaktoryPublishService) {
	fs = middleware.PrometheusPublishTaskMiddleware{Next: fs}
	fs = middleware.StatusRequestMiddleware{StatusClient: statusClient, Next: fs}
	requestHandler := httptransport.NewServer(
		transport.CreateRestHandler(fs),
		transport.DecodeScreenshotTask,
		transport.EncodeResponse,
	)
	http.Handle("/", requestHandler)
	fmt.Println(http.ListenAndServe(":"+os.Getenv("OPTIMIZATION_SERVICE_EXTERNAL_PORT"), nil))
}
