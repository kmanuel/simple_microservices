package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type TaskType string

type Task struct {
	ID   string   `json:"id"`
	Type TaskType `json:"tasktype"`
}

const (
	OptimizeImage               = "optimize_image"
	CropImage                   = "crop_image"
	FaceDetection               = "face_detection"
	FullPageScreenshot          = "full_page_screenshot"
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

func NewTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = uuid.New().String()
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(tasks)

	dispatchTask(&task)
	//publishToFactory(&task)
}

func dispatchTask(t *Task) {
	switch t.Type {
	case FaceDetection:
		var jsonBody = []byte(`{ 
			"in": "test_image.jpg", 
			"out": "test_output.jpg", 
			"width": 10, 
			"height": 10
		}`)
		http.Post("http://localhost:8081/crop", "application/json", bytes.NewBuffer(jsonBody))
		break
	case CropImage:
		var jsonBody = []byte(`{
			"in": "test.jpg", 
			"out": "res.jpg",
			"width": 20,
			"height": 20
		}`)
		http.Post("http://localhost:8082/crop", "application/json", bytes.NewBuffer(jsonBody))
		break
	}
}

//func publishToFactory(t *Task) {
//	client, err := faktory.Open()
//	log.Println(err)
//	job := faktory.NewJob("SomeJob", &t)
//	err = client.Push(job)
//	log.Println(err)
//	log.Println("published task to factory")
//}
