package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/resolver"
	"github.com/manyminds/api2go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type NewTask struct {
	TaskId string `json:"id"`
}

type TaskStatus struct {
	gorm.Model
	TaskId string `json:"task_id"`
	Status string `json:"status"`
}

type TaskStatusUpdate struct {
	Status string `json:"status"`
}

var dbHost string
var dbPort string
var dbUser string
var dbName string
var dbPassword string

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dbHost = os.Getenv("POSTGRES_HOST")
	dbPort = os.Getenv("POSTGRES_PORT")
	dbUser = os.Getenv("POSTGRES_USER")
	dbName = os.Getenv("POSTGRES_DB")
	dbPassword = os.Getenv("POSTGRES_PASSWORD")

	db, err := openDb()
	defer db.Close()
	if err != nil {
		panic(err)
		return
	}
	db.AutoMigrate(&TaskStatus{})

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
	port := 8080
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port})
	handler := api.Handler().(*httprouter.Router)
	handler.POST("/tasks/", CreateNew)
	handler.GET("/tasks", GetTasks)
	handler.POST("/tasks/:taskId/status", UpdateStatus)
	handler.GET("/tasks/:taskId/status", GetStatus)
	handler.GET("/health", GetHealth)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func CreateNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("creating new TaskStatus")

	var newTask NewTask
	_ = json.NewDecoder(r.Body).Decode(&newTask)

	taskStatus := TaskStatus{
		TaskId: newTask.TaskId,
		Status: "new",
	}

	db, err := openDb()
	defer db.Close()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	db.Create(&taskStatus)

	w.WriteHeader(201)
}

func GetTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("getting all tasks")

	db, err := openDb()
	defer db.Close()
	if err != nil {
		log.Error("failed to open db", err)
		w.WriteHeader(500)
		return
	}

	var tasks []TaskStatus
	if err := db.Find(&tasks).Error; err != nil {
		log.Error("failed to fetch all taskStatus from db")
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Error("failed to write response")
	}
}

func UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Info("received update request")

	var statusUpdate TaskStatusUpdate
	_ = json.NewDecoder(r.Body).Decode(&statusUpdate)

	db, err := openDb()
	defer db.Close()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var taskStatus TaskStatus
	if err := db.Where("task_id = ? ", ps.ByName("taskId")).First(&taskStatus).Error; err != nil {
		w.WriteHeader(500)
		return
	}

	taskStatus.Status = statusUpdate.Status
	db.Save(&taskStatus)

	w.WriteHeader(200)
}

func GetStatus(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	log.Info("received GET status request")

	db, err := openDb()
	defer db.Close()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var taskStatus TaskStatus
	if err := db.Where("task_id = ? ", ps.ByName("taskId")).First(&taskStatus).Error; err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(taskStatus)
	if err != nil {
		log.Error("failed to write response")
	}
}

func GetHealth(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
}

func openDb() (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open("postgres", psqlInfo)
	return db, err
}
