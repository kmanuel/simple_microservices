package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	faktory "github.com/contribsys/faktory/client"
)

type TaskType string

type Task struct {
	ID string `json:"id"`
	Type TaskType `json:"tasktype"`
}

const (
	OptimizeImage = "optimize_image"
	CropImage = "crop_image"
	FaceDetection = "face_detection"
	FullPageScreenshot = "full_page_screenshot"
	ExtractMostSignificantImage = "extract_most_significant_image"
)

var tasks []Task

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks", NewTask).Methods("POST")
	log.Println(http.ListenAndServe(":8080", router))
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func NewTask(w http.ResponseWriter, r * http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = uuid.New().String()
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(tasks)

	publishToFactory(&task)
}

func publishToFactory(t *Task) {
	client, err := faktory.Open()
	log.Println(err)
	job := faktory.NewJob("SomeJob", &t)
	err = client.Push(job)
	log.Println(err)
	log.Println("published task to factory")
}
