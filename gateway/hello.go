package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/google/uuid"
)


type TaskType string

const (
	OptimizeImage = "optimize_image"
	CropImage = "crop_image"
	FaceDetection = "face_detection"
	FullPageScreenshot = "full_page_screenshot"
	ExtractMostSignificantImage = "extract_most_significant_image"
)


func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks", NewTask).Methods("POST")
	log.Println(http.ListenAndServe(":8080", router))
}

type Task struct {
	ID string `json:"id"`
	Type TaskType `json:"tasktype"`
}

var tasks []Task

func GetTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func NewTask(w http.ResponseWriter, r * http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = uuid.New().String()
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(tasks)
}
