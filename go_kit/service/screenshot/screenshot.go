package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/screenshot/status_client"
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

	var statusClient status_client.StatusClient
	statusClient = status_client.NewStatusClient(taskType)

	var screenshotService service.ScreenshotService
	screenshotService = service.NewScreenshotService()

	var faktoryService service.FaktoryService
	faktoryService = service.NewFaktoryService(taskType)

	go startPrometheus()

	var faktoryListenService service.FaktoryListenService
	faktoryListenService = faktoryService
	go startFaktory(faktoryListenService, screenshotService, statusClient)

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
	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}

func startFaktory(fs service.FaktoryListenService, screenshotService service.ScreenshotService, statusClient status_client.StatusClient) {
	screenshotService = middleware.StatusPerformMiddleware{StatusClient: statusClient, Next: screenshotService}
	screenshotService = middleware.PrometheusProcessTaskMiddleware{Next: screenshotService}
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(screenshotService))
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
	fmt.Println(http.ListenAndServe(":8080", nil))
}
