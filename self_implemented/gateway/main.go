package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/gateway/resolver"
	"github.com/manyminds/api2go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Task struct {
	ID         string
	Type       string                 `json:"tasktype"`
	TaskParams map[string]interface{} `json:"taskParams"`
}

type NewTaskType struct {
	Id string `json:"id"`
}

type UploadResponse struct {
	FileId string `json:"fileId"`
}

var tasks []Task

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))

	go startPrometheus()

	startRestApi()
}

func startRestApi() {
	port := 8080
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port})
	handler := api.Handler().(*httprouter.Router)
	handler.GET("/tasks", GetTasks)
	handler.POST("/tasks", NewTask)
	handler.POST("/upload", UploadFile)
	handler.GET("/tasks/:taskId/download", DownloadFile)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"service", "type"},
	)
)

func startPrometheus() {
	prometheus.MustRegister(requests)

	var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}


func UploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requests.With(prometheus.Labels{"service":"gateway", "type": "upload"}).Inc()
	log.Info("incoming file upload request")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		panic(err)
	}

	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	uploadedFileName := uuid.New().String()
	minioconnector.UploadFileStream(file, uploadedFileName)

	var uploadResponse UploadResponse
	uploadResponse.FileId = uploadedFileName

	err = json.NewEncoder(w).Encode(uploadResponse)
	if err != nil {
		log.Error("error writing response")
	}
}

func DownloadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requests.With(prometheus.Labels{"service":"gateway", "type": "download"}).Inc()
	taskId := ps.ByName("taskId")
	log.Info("download request for taskId=", taskId)
	object := minioconnector.GetObject(taskId)
	io.Copy(w, object)
}

func GetTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requests.With(prometheus.Labels{"service":"gateway", "type": "get_tasks"}).Inc()
	log.Info("received request for all tasks")

	requestServiceUrl, e := url.Parse("http://request_service:8080")
	if e != nil {
		panic(e)
	}
	httputil.NewSingleHostReverseProxy(requestServiceUrl).ServeHTTP(w, r)

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Error("error writing response")
	}

}

func NewTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requests.With(prometheus.Labels{"service":"gateway", "type": "create_task"}).Inc()
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
	_ = t.UnmarshalJSON(body)

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

func (t *Task) UnmarshalJSON(data []byte) error {
	var jsonMap map[string]interface{}

	if t == nil {
		return errors.New("RawString: UnmarshalJSON on nil pointer")
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
