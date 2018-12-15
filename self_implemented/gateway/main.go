package main

import (
	"bytes"
	"encoding/json"
	"errors"
	faktory "github.com/contribsys/faktory/client"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type Task struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"tasktype"`
	TaskParams map[string]interface{} `json:"taskParams"`
}

const (
	OptimizeImage               = "optimize"
	CropImage                   = "crop"
	FaceDetection               = "face_detection"
	Screenshot                  = "screenshot"
	ExtractMostSignificantImage = "most_significant_image"
)

var tasks []Task

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks", NewTask).Methods("POST")
	log.Info(http.ListenAndServe(":8080", router))
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func NewTask(w http.ResponseWriter, r *http.Request) {
	log.Info("received request for new task")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var t Task
	_ = t.UnmarshalJSON(body)

	dispatchTask(&t)

	log.WithFields(log.Fields{
		"taskID": t.ID,
	}).Info("finished task handling")
	publishToFactory(&t)
}

func dispatchTask(t *Task) {
	log.WithFields(log.Fields{
		"taskID": t.ID,
		"type": t.Type,
	}).Info("dispatching new Task")

	switch t.Type {
	case FaceDetection:
		taskParameters := ExtractTaskParameters(t)
		http.Post("portrait:8080", "application/json", bytes.NewBuffer([]byte(taskParameters)))
		break
	case CropImage:
		taskParameters := ExtractTaskParameters(t)
		http.Post("http://crop:8080", "application/json", bytes.NewBuffer([]byte(taskParameters)))
		break
	case Screenshot:
		taskParameters := ExtractTaskParameters(t)
		http.Post("http://screenshot:8080", "application/json", bytes.NewBuffer([]byte(taskParameters)))
		break
	case OptimizeImage:
		taskParameters := ExtractTaskParameters(t)
		http.Post("http://optimize:8080", "application/json", bytes.NewBuffer([]byte(taskParameters)))
		break
	case ExtractMostSignificantImage:
		taskParameters := ExtractTaskParameters(t)
		http.Post("http://most_significant_image:8080", "application/json", bytes.NewBuffer([]byte(taskParameters)))
		break
	default:
		log.WithFields(log.Fields{
		"taskID": t.ID,
		}).Warn("no handler found for task")
	}
	log.WithFields(log.Fields{
		"taskID": t.ID,
	}).Debug("Dispatching of task finished for task with ID=" + t.ID)
}

func ExtractTaskParameters(t *Task) string {
	log.WithFields(log.Fields{
		"taskID": t.ID,
	}).Debug("extracting parameters from task")

	var params []string
	for key, val := range t.TaskParams {
		if key != "id" && key != "tasktype" {
			t.TaskParams[key] = val
			param := `"` + key + `": "` + `"` + val.(string) + `"`
			params = append(params, param)
		}
	}
	jsonString := `{` + strings.Join(params, ", \n") + `}`

	log.WithFields(log.Fields{
		"taskID": t.ID,
	}).Debug("extracted task parameters" + jsonString + " from task with ID=" + t.ID)

	return jsonString
}

func (t *Task) UnmarshalJSON(data []byte) error {
	var jsonMap map[string]interface{}

	if t == nil {
		return errors.New("RawString: UnmarshalJSON on nil pointer")
	}

	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return err
	}

	t.ID = jsonMap["id"].(string)
	t.Type = jsonMap["tasktype"].(string)

	t.TaskParams = make(map[string]interface{})

	for key, val := range jsonMap {
		if key != "id" && key != "tasktype" {
			t.TaskParams[key] = val
		}
	}

	return nil
}

func publishToFactory(t *Task) {
	client, err := faktory.Open()
	log.Println(err)
	job := faktory.NewJob(t.Type, &t.TaskParams)
	job.Queue = t.Type
	job.Custom = t.TaskParams
	err = client.Push(job)
	log.Println(err)
	log.Println("published task to factory")
}
