package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/portrait/controller"
	"github.com/kmanuel/simple_microservices/self_implemented/service/portrait/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)


var taskType = "portrait"

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_handle_count",
			Help: "Number of handled requests.",
		},
		[]string{"type"},
	)
)

func main() {
	initMinio()

	var taskStatusService service.TaskStatusService
	taskStatusService = service.NewTaskStatusService()

	var taskService service.TaskService
	taskService = service.NewTaskService(requests, taskType)

	go startPrometheus()
	go startFaktoryListener(taskStatusService, taskService)
	startRestApi(taskStatusService)
}

func initMinio() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
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
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func startFaktoryListener(taskStatusService service.TaskStatusService, taskService service.TaskService) {
	service.NewFactoryListenerService(taskStatusService, taskService, taskType)
}

func startRestApi(taskStatusService service.TaskStatusService) {
	var faktoryPublishService service.FaktoryPublishService
	faktoryPublishService = service.NewFaktoryPublishService(taskType)

	var taskHandler handler.TaskHandler
	taskHandler = handler.NewTaskHandler(faktoryPublishService, taskStatusService, taskType)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/" + taskType, taskHandler.PerformTask).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}
