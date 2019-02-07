package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
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
	dispatchCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dispatch_count",
			Help: "Number of dispatchCounter handled from faktory.",
		},
		[]string{"type"},
	)
)

func main() {
	log.Info("starting gateway")

	loadDotEnv()
	initMinio()
	go startPrometheus()
	startJsonRestApi()
}

func loadDotEnv() {
	log.Debug("Loading dotenv")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func initMinio() {
	log.Debug("initializing minio")

	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
}

func startPrometheus() {
	log.Debug("starting prometheus")

	prometheus.MustRegister(dispatchCounter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func startJsonRestApi() {
	log.Debug("starting REST API")


	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")
	log.Info("UPDATED")

	rootHandler := &RootHandler{}
	imageHandler := &api_image.ImageHandler{DispatchCounter: dispatchCounter}
	proxyHandler := &api_task.ProxyHandler{DispatchCounter: dispatchCounter}
	imageTransformationHandler := &api_image.ImageTaskHandler{DispatchCounter: dispatchCounter}

	myRouter := mux.NewRouter().StrictSlash(false)

	myRouter.HandleFunc("/", rootHandler.GetRootResource).Methods("GET")

	myRouter.HandleFunc("/images", imageHandler.UploadImage).Methods("POST")
	myRouter.HandleFunc("/images/{id}", imageHandler.DownloadImage).Methods("GET")
	myRouter.HandleFunc("/images/{id}/tasks", imageTransformationHandler.HandleGetTasks).Methods("GET")
	myRouter.HandleFunc("/tasks", proxyHandler.ProxyToRequestService)
	myRouter.HandleFunc("/screenshot", proxyHandler.CreateScreenshotTask)
	myRouter.HandleFunc("/crop", proxyHandler.CreateCropTask)
	myRouter.HandleFunc("/most_significant_image", proxyHandler.CreateMostSignificantImageTask)
	myRouter.HandleFunc("/optimization", proxyHandler.CreateOptimizationTask)
	myRouter.HandleFunc("/portrait", proxyHandler.CreatePortraitTask)
	myRouter.HandleFunc("/info", proxyHandler.GetFaktoryInfo)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
