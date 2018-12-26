package main

import (
	"fmt"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_faktory"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_root"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_image"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/api/api_task"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	//"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
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
	jsonapi.Instrumentation = func(r *jsonapi.Runtime, eventType jsonapi.Event, callGUID string, dur time.Duration) {
		metricPrefix := r.Value("instrument").(string)

		if eventType == jsonapi.UnmarshalStart {
			fmt.Printf("%s: id, %s, started at %v\n", metricPrefix+".jsonapi_unmarshal_time", callGUID, time.Now())
		}

		if eventType == jsonapi.UnmarshalStop {
			fmt.Printf("%s: id, %s, stopped at, %v , and took %v to unmarshal payload\n", metricPrefix+".jsonapi_unmarshal_time", callGUID, time.Now(), dur)
		}

		if eventType == jsonapi.MarshalStart {
			fmt.Printf("%s: id, %s, started at %v\n", metricPrefix+".jsonapi_marshal_time", callGUID, time.Now())
		}

		if eventType == jsonapi.MarshalStop {
			fmt.Printf("%s: id, %s, stopped at, %v , and took %v to marshal payload\n", metricPrefix+".jsonapi_marshal_time", callGUID, time.Now(), dur)
		}
	}

	myRouter := mux.NewRouter().StrictSlash(false)

	rootHandler := &api_root.RootHandler{}
	imageHandler := &api_image.ImageHandler{ RequestCounter: requests }
	taskHandler := &api_task.TaskHandler{ RequestCounter: requests }
	faktoryHandler := &api_faktory.FaktoryHandler{RequestCounter: requests}

	myRouter.HandleFunc("/resources", rootHandler.ServeHTTP)
	myRouter.HandleFunc("/upload", imageHandler.ServeUploadHTTP)
	myRouter.HandleFunc("/image/{id}", imageHandler.ServeHTTP)
	myRouter.HandleFunc("/tasks", taskHandler.ServeHTTP)
	myRouter.HandleFunc("/faktory/info", faktoryHandler.ServeHTTP)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
