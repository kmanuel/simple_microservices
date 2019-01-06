package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/crop/status_client"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	initMinio()

	var cropService service.CropService
	cropService = service.NewCropService()

	var statusClient status_client.StatusClient
	statusClient = status_client.NewStatusClient()

	var faktoryService service.FaktoryService
	faktoryService = service.NewFaktoryService(cropService)

	go startPrometheus()
	go startFaktory(faktoryService, cropService, statusClient)
	startExternalApi(statusClient, faktoryService)
}

func initMinio() {
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startFaktory(fs service.FaktoryListenService, cropService service.CropService, statusClient status_client.StatusClient) {
	cropService = middleware.StatusCropMiddleware{StatusClient: statusClient, Next: cropService}
	handler := transport.CreateFaktoryHandler(cropService)
	fs.Handle("crop", handler)
}

func startExternalApi(statusClient status_client.StatusClient, fs service.FaktoryPublishService) {
	fs = middleware.PrometheusPublishTaskMiddleware{Next: fs}
	fs = middleware.RequestStatusMiddleware{StatusClient: statusClient, Next: fs}
	requestHandler := httptransport.NewServer(
		transport.MakeCropRequestHandler(fs),
		transport.DecodeCropTask,
		transport.EncodeResponse,
	)
	http.Handle("/", requestHandler)
	fmt.Println(http.ListenAndServe(":"+os.Getenv("CROP_SERVICE_EXTERNAL_PORT"), nil))
}

func startPrometheus() {
	prometheus.MustRegister(requests)
	var addr = flag.String("listen-address", ":"+os.Getenv("CROP_SERVICE_PROMETHEUS_PORT"), "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}
