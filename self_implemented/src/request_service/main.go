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
	myRouter.HandleFunc("/tasks", taskHandler.ServeHTTP)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), handler))
	//
	//api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: 8080})
	//handler := api.Handler().(*httprouter.Router)
	//
	//handler.POST("/tasks/:taskId/status", UpdateStatus)
	//handler.GET("/tasks/:taskId/status", GetStatus)
	//handler.GET("/health", GetHealth)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), handler))
}

//
//
//func UpdateStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	log.Info("received update request")
//
//	var statusUpdate TaskStatusUpdateDTO
//	_ = json.NewDecoder(r.Body).Decode(&statusUpdate)
//
//	db, err := openDb()
//	defer db.Close()
//	if err != nil {
//		w.WriteHeader(500)
//		return
//	}
//
//	var taskStatus TaskStatus
//	if err := db.Where("task_id = ? ", ps.ByName("taskId")).First(&taskStatus).Error; err != nil {
//		w.WriteHeader(500)
//		return
//	}
//
//	taskStatus.Status = statusUpdate.Status
//	db.Save(&taskStatus)
//
//	w.WriteHeader(200)
//}
//
//func GetStatus(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
//	log.Info("received GET status request")
//
//	db, err := openDb()
//	defer db.Close()
//	if err != nil {
//		w.WriteHeader(500)
//		return
//	}
//
//	var taskStatus TaskStatus
//	if err := db.Where("task_id = ? ", ps.ByName("taskId")).First(&taskStatus).Error; err != nil {
//		w.WriteHeader(500)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(taskStatus)
//	if err != nil {
//		log.Error("failed to write response")
//	}
//}
//
//func GetHealth(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
//	w.WriteHeader(200)
//}
