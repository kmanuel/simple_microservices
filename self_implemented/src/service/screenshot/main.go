package main

import (
	"encoding/json"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/screenshot/update_status"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
)

type Request struct {
	Url string `json:"url"`
}

type TaskParamsType map[string]interface{}

type Task struct {
	ID         string         `json:"id"`
	Type       string         `json:"tasktype"`
	TaskParams TaskParamsType `json:"taskParams"`
}

func main() {
	godotenv.Load()
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))

	//router := mux.NewRouter()
	//router.HandleFunc("/", HandleRequest).Methods("POST")
	//log.Info(http.ListenAndServe(":8080", router))

	startFaktory()
}

func startFaktory() {
	mgr := worker.NewManager()
	mgr.Use(func(perform worker.Handler) worker.Handler {
		return func(ctx worker.Context, job *faktory.Job) error {
			log.Printf("Starting work on job %s of type %s with custom %v\n", ctx.Jid(), ctx.JobType(), job.Custom)
			err := perform(ctx, job)
			log.Printf("Finished work on job %s with error %v\n", ctx.Jid(), err)
			return err
		}
	})
	mgr.Register("screenshot", convertTask)
	mgr.Queues = []string{"screenshot"}
	var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
	})
	// Start processing jobs, this method does not return
	mgr.Run()
}

func convertTask(ctx worker.Context, args ...interface{}) error {
	log.Info("Working on job %s\n", ctx.Jid())


	strings, ok := args[0].(map[string]interface{})
	if !ok {
		log.Error("couldnt convert args[0]")
	} else {
		log.Error("url: ", strings["url"])

		update_status.NotifyAboutProcessingStart(strings["id"].(string))

		takeScreenShot(strings["url"].(string))

		update_status.NotifyAboutCompletion(strings["id"].(string))
	}

	return nil
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	log.Info("received request")
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	outputFilePath := takeScreenShot(task.Url)

	minioconnector.UploadFile(outputFilePath)

	log.Info("finished request")
}

func takeScreenShot(url string) string {
	log.WithField("url", url).Info("taking screenshot")

	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "./lib/bin/phantomjs"
	jsPath := "./lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		log.Printf("process job err - %s\n", err.Error())
	}

	log.WithField("url", url).Info("took screenshot")

	return outputFilePath
}
