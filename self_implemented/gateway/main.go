package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/controller/image"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/controller/task"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/resolver"
	"github.com/manyminds/api2go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"controller", "type"},
	)
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	initMinio()
	go startPrometheus()
	startRestApi()
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
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func startRestApi() {
	imageController := image.NewImageController(requests)
	taskController := task.NewTaskController(requests)

	port := 8080
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port})
	handler := api.Handler().(*httprouter.Router)

	handler.GET("/tasks", taskController.HandleGetTasks)
	handler.POST("/tasks", taskController.HandleTaskCreation)
	handler.GET("/faktory/info", taskController.HandleGetTasksInfo)

	handler.POST("/upload", imageController.HandleUpload)
	handler.GET("/tasks/:taskId/download", imageController.HandleDownload)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
