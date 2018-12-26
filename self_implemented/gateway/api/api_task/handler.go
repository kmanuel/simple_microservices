package api_task

import (
	"bytes"
	"encoding/json"
	"errors"
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type NewTaskType struct {
	Id string `json:"id"`
}

type TaskHandler struct{
	RequestCounter *prometheus.CounterVec
}

type Task struct {
	ID         string
	Type       string                 `json:"tasktype"`
	TaskParams map[string]interface{} `json:"taskParams"`
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc
	switch r.Method {
	case http.MethodPost:
		methodHandler = h.HandleTaskCreation
	case http.MethodGet:
		methodHandler = h.HandleGetTasks
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	methodHandler(w, r)
}

func (h *TaskHandler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	h.RequestCounter.With(prometheus.Labels{"controller":"gateway", "type": "get_tasks"}).Inc()
	log.Info("received request for all tasks")

	requestServiceUrl, e := url.Parse("http://request_service:8080")
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)
}


func (h *TaskHandler) HandleTaskCreation(w http.ResponseWriter, r *http.Request) {
	h.RequestCounter.With(prometheus.Labels{"controller":"gateway", "type": "create_task"}).Inc()
	log.Info("received request for new task")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	taskId := uuid.New().String()

	err = sendToRequestService(taskId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var t Task
	_ = t.unmarshalJSON(body)

	t.ID = taskId

	log.WithFields(log.Fields{
	}).Info("finished task handling")
	err = publishToFactory(&t)
	if err != nil {
		w.WriteHeader(500)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		err := json.NewEncoder(w).Encode(t)
		if err != nil {
			log.Error("error writing response")
		}
	}
}

func sendToRequestService(taskId string) error {
	var nt NewTaskType
	nt.Id = taskId
	marshal, e := json.Marshal(nt)
	if e != nil {
		panic(e)
	}
	_, err := http.Post("http://request_service:8080/tasks", "application/json", bytes.NewBuffer([]byte(marshal)))
	return err
}

func publishToFactory(t *Task) error {
	log.Info("publish to faktory")
	client, err := faktory.Open()
	if err != nil {
		return err
	}
	job := faktory.NewJob(t.Type, &t.TaskParams)
	job.Queue = t.Type
	t.TaskParams["id"] = t.ID
	job.Custom = t.TaskParams
	err = client.Push(job)

	return err
}


func (t *Task) unmarshalJSON(data []byte) error {
	var jsonMap map[string]interface{}

	if t == nil {
		return errors.New("RawString: unmarshalJSON on nil pointer")
	}

	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return err
	}

	t.Type = jsonMap["tasktype"].(string)

	t.TaskParams = make(map[string]interface{})

	for key, val := range jsonMap {
		if key != "id" && key != "tasktype" {
			t.TaskParams[key] = val
		}
	}

	return nil
}