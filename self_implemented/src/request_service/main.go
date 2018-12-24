package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/database"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/resolver"
	"github.com/manyminds/api2go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type NewTask struct {
	Id	string	`json:"id"`
}


type TaskStatus struct {
	Id 		string	`json:"id"`
	Status 	string 	`json:"status"`
}

type TaskStatusUpdate struct {
	Status	string	`json:"status"`
}

func main() {
	godotenv.Load()
	dbPortStr := os.Getenv("POSTGRES_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		panic(err)
	}

	database.Init(
		os.Getenv("POSTGRES_HOST"),
		dbPort,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	go startPrometheus()

	startRestApi()
}

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"service", "status"},
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
	port := 8080
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port})
	handler := api.Handler().(*httprouter.Router)
	handler.POST("/tasks/", CreateNew)
	handler.GET("/tasks", GetTasks)
	handler.POST("/tasks/:taskId/status", UpdateStatus)
	handler.GET("/tasks/:taskId/status", GetStatus)
	handler.GET("/health", GetHealth)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func CreateNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("creating new TaskStatus")

	var newTask NewTask
	_ = json.NewDecoder(r.Body).Decode(&newTask)

	database.Persist(newTask.Id)
}

func GetTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("getting all tasks")

	all := database.FetchAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Info("received update request")

	var statusUpdate TaskStatusUpdate
	_ = json.NewDecoder(r.Body).Decode(&statusUpdate)

	taskId := ps.ByName("taskId")

	database.UpdateStatus(taskId, statusUpdate.Status)


	log.Error("task with id="+taskId+" gets updated status=" + statusUpdate.Status)
}

func GetStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Info("received GET status request")

	taskId := ps.ByName("taskId")

	status := database.FetchStatus(taskId)

	var t TaskStatus
	t.Id = taskId
	t.Status = status

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func GetHealth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
}
