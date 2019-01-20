package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_faktory"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_image"
	. "github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_root"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_task"
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
	log.Info("starting gateway")

	log.Debug("Loading dotenv")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log.Debug("initializing minio")
	initMinio()

	log.Debug("starting prometheus")
	go startPrometheus()

	log.Debug("starting REST API")
	startJsonRestApi()
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

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func startJsonRestApi() {
	myRouter := mux.NewRouter().StrictSlash(false)

	rootHandler := &RootHandler{}
	imageHandler := &api_image.ImageHandler{RequestCounter: requests}
	taskHandler := &api_task.TaskHandler{RequestCounter: requests}
	faktoryHandler := &api_faktory.FaktoryHandler{RequestCounter: requests}
	imageTransformationHandler := &api_image.ImageTaskHandler{RequestCounter: requests}

	myRouter.HandleFunc("/", rootHandler.ServeHTTP)
	myRouter.HandleFunc("/images", imageHandler.ServeUploadHTTP)
	myRouter.HandleFunc("/images/{id}", imageHandler.ServeDownload)
	myRouter.HandleFunc("/images/{id}/tasks", imageTransformationHandler.ServeHTTP)
	myRouter.HandleFunc("/tasks", taskHandler.ServeHTTP)
	myRouter.HandleFunc("/screenshot", taskHandler.ServeScreenshotHTTP)
	myRouter.HandleFunc("/crop", taskHandler.ServeCropHTTP)
	myRouter.HandleFunc("/most_significant_image", taskHandler.ServeMostSignificantHTTP)
	myRouter.HandleFunc("/optimization", taskHandler.ServeOptimizationHTTP)
	myRouter.HandleFunc("/portrait", taskHandler.ServePortraitHTTP)
	myRouter.HandleFunc("/faktory/info", faktoryHandler.ServeHTTP)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
