package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/joho/godotenv"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/middleware"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/service"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/status_client"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/transport"
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

const taskType = "most_significant_image"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var statusClient status_client.StatusClient

	statusClient = status_client.NewStatusClient(taskType)

	var faktoryService service.FaktoryService
	faktoryService = service.NewFaktoryService(taskType)

	var imageService service.MostSignificantImageService
	imageService = service.NewMostSignificantImageService()

	go startPrometheus()
	go startFaktory(faktoryService, imageService, statusClient)
	startRestApi(statusClient, &faktoryService)
}

func startRestApi(sc status_client.StatusClient, s *service.FaktoryService) {
	var publishService service.FaktoryPublishService
	publishService = middleware.StatusRequestMiddleware{StatusClient: sc, Next: publishService}
	publishService = middleware.PrometheusPublishTaskMiddleware{Next: *s}

	requestHandler := httptransport.NewServer(
		transport.CreateRestHandler(*s),
		transport.DecodeMostSignificantImageTask,
		transport.EncodeResponse,
	)
	http.Handle("/", requestHandler)
	fmt.Println(http.ListenAndServe(":"+os.Getenv("MOST_SIGNIFICANT_IMAGE_EXTERNAL_PORT"), nil))
}

func startFaktory(fs service.FaktoryListenService, s service.MostSignificantImageService, sc status_client.StatusClient) {
	s = middleware.StatusPerformMiddleware{StatusClient: sc, Next: s}
	s = middleware.PrometheusProcessTaskMiddleware{Next: s}
	fs.Handle(taskType, transport.CreateFaktoryListenHandler(s))
}

func startPrometheus() {
	prometheus.MustRegister(requests)
	var addr = flag.String("listen-address", ":"+os.Getenv("CROP_SERVICE_PROMETHEUS_PORT"), "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(http.ListenAndServe(*addr, nil))
}

