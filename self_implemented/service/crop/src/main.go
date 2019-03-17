package main

import (
	"flag"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/service/crop/src/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var taskType = "crop"

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
	log.Info("starting crop service")

	log.Info("inputBucketName=" + os.Getenv("INPUT_BUCKET_NAME"))

	minioService := minioconnector.NewMinioService(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("INPUT_BUCKET_NAME"),
		taskType)

	go startPrometheus()
	startFaktoryListener(minioService)
}

func startPrometheus() {
	prometheus.MustRegister(requests)
	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func startFaktoryListener(minioConnector *minioconnector.MinioService) {
	taskService := service.NewTaskService(minioConnector, requests, taskType)

	faktoryService := service.NewFaktoryListenerService(taskService, taskType)

	err := faktoryService.Start()
	if err != nil {
		panic(err)
	}
}
