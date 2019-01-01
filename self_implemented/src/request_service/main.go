package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var dbHost string
var dbPort string
var dbUser string
var dbName string
var dbPassword string

func main() {
	api.InitDb()

	go startPrometheus()

	startRestApi()
}

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"controller", "status"},
	)
)

func startPrometheus() {
	prometheus.MustRegister(requests)

	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func startRestApi() {
	myRouter := mux.NewRouter().StrictSlash(false)

	taskHandler := &api.TaskHandler{RequestCounter: requests}

	myRouter.
		Path("/tasks").
		Methods(http.MethodGet).
		HandlerFunc(taskHandler.GetTasks)
	myRouter.
		Path("/tasks").
		Methods(http.MethodPost).
		HandlerFunc(taskHandler.CreateTask)

	myRouter.
		Path("/tasks/status/{id}").
		Methods(http.MethodPost).
		HandlerFunc(taskHandler.UpdateTask)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
