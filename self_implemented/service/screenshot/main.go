package main

import (
	"flag"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/screenshot/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const taskType = "screenshot"

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
	go startPrometheus()
	startFaktoryListener()
}

func initMinio() {
	hystrix.ConfigureCommand("update_task_status", hystrix.CommandConfig{
		Timeout:               60000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

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

func startFaktoryListener() {
	var statusService service.TaskStatusService
	statusService = service.NewTaskStatusService()

	var taskService service.TaskService
	taskService = service.NewTaskService(requests, taskType)

	faktoryService := service.NewFaktoryListenerService(statusService, taskService, taskType)

	err := faktoryService.Start()
	if err != nil {
		panic(err)
	}
}
