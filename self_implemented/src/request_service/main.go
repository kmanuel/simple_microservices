package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/database"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/resolver"
	"github.com/manyminds/api2go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

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

	status := database.FetchStatus("otherqwer")

	log.Error("status of otherqwer task is" + status)


	port := 8080
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port})

	handler := api.Handler().(*httprouter.Router)
	handler.POST("/tasks/{taskId}/status", UpdateStatus)

	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)

	//tasks := database.FetchTasks()
	//
	//for _, task := range *tasks {
	//	log.Error("taskId=" + task.ID)
	//}

}

func UpdateStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("received request for new task")

	params := mux.Vars(r)

	var statusUpdate TaskStatusUpdate
	_ = json.NewDecoder(r.Body).Decode(&statusUpdate)

	taskId := params["taskId"]


	log.Error("task with id="+taskId+" gets updated status=" + statusUpdate.Status)



	//
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//var t model.Task
	//_ = t.UnmarshalJSON(body)
	//
	//t.ID = uuid.New().String()
	//
	//log.WithFields(log.Fields{
	//	"taskID": t.ID,
	//}).Info("finished task handling")
	//publishToFactory(&t)
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(201)
	//json.NewEncoder(w).Encode(t)
}

