package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
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

	var cropService service.ImageService
	cropService = service.NewCropService()

	var statusClient service.StatusClient
	statusClient = service.NewStatusClient()

	go startPrometheus()
	startExternalApi(cropService, statusClient)
}

func initMinio() {
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startExternalApi(imageService service.ImageService, statusClient service.StatusClient) {
	imageService = middleware.NewPrometheusMiddleware(imageService, taskType)
	imageService = middleware.NewRequestStatusMiddleware(statusClient, imageService)

	requestHandler := httptransport.NewServer(
		transport.MakeTaskHandler(imageService),
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
