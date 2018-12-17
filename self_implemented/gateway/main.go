package main

import (
	"encoding/json"
	"errors"
	"fmt"
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/kmanuel/minioconnector"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/examples/resolver"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Task struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"tasktype"`
	TaskParams map[string]interface{} `json:"taskParams"`
}

type UploadResponse struct {
	FileId	string	`json:"fileId"`
}

var tasks []Task

func main() {
	godotenv.Load()
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))

	port := 8080
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port})

	handler := api.Handler().(*httprouter.Router)
	handler.GET("/tasks", GetTasks)
	handler.POST("/tasks", NewTask)
	handler.POST("/upload", UploadFile)

	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func UploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("incoming file upload request")
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("/tmp/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	uploadedFileName := minioconnector.UploadFile("/tmp/" + handler.Filename)

	var uploadResponse UploadResponse
	uploadResponse.FileId = uploadedFileName

	json.NewEncoder(w).Encode(uploadResponse)

}

func GetTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(tasks)
}

func NewTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("received request for new task")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var t Task
	_ = t.UnmarshalJSON(body)

	t.ID = uuid.New().String()

	log.WithFields(log.Fields{
		"taskID": t.ID,
	}).Info("finished task handling")
	publishToFactory(&t)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(t)
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
