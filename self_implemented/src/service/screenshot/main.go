package main

import (
	"bytes"
	"flag"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/screenshot/update_status"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
)

type Request struct {
	Url string `json:"url"`
}

type Task struct {
	ID  string `jsonapi:"primary,screenshot_task"`
	Url string `jsonapi:"attr,url"`
}

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

	startFaktory()
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

	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
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

	task := new(Task)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		log.Error("failed to deserialize task", args)
		return err
	}

	update_status.NotifyAboutProcessingStart(task.ID)

	outputFilePath, err := takeScreenShot(task.Url)
	if err != nil {
		_ = ctx.Err()
		return nil
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, task.ID)
	if err != nil {
		_ = ctx.Err()
		return nil
	}

	update_status.NotifyAboutCompletion(task.ID)
	ctx.Done()

	return nil
}

func takeScreenShot(url string) (string, error) {
	log.WithField("url", url).Info("taking screenshot")

	chromeUserAgent := "Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	phantomJSBin := "./lib/bin/phantomjs"
	jsPath := "./lib/js/screenshot.js"
	logFile := "screenshot_service.log"

	outputFilePath := "/tmp/output" + uuid.New().String() + ".jpg"

	cmd := exec.Command(phantomJSBin, jsPath, url, outputFilePath, logFile, chromeUserAgent)

	if err := cmd.Run(); nil != err {
		return "", err
	}

	log.WithField("url", url).Info("took screenshot")

	return outputFilePath, nil
}
