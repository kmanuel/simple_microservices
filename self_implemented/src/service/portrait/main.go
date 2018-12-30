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
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/portrait/update_status"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type Task struct {
	ID			string		`jsonapi:"primary,portrait_task"`
	ImageId		string		`jsonapi:"attr,image_id"`
	Width 		int			`jsonapi:"attr,width"`
	Height 		int			`jsonapi:"attr,height"`
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
	mgr.Register("portrait", convertTask)
	mgr.Queues = []string{"portrait"}
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

	downloadedFilePath, err := minioconnector.DownloadFile(task.ImageId)
	if err != nil {
		_ = ctx.Err()
		return err
	}

	outputFilePath, err := ExtractPortrait(downloadedFilePath, task.Width, task.Height)
	if err != nil {
		_ = ctx.Err()
		return err
	}

	_, err = minioconnector.UploadFileWithName(outputFilePath, task.ID)
	if err != nil {
		_ = ctx.Err()
		return err
	}

	update_status.NotifyAboutCompletion(task.ID)

	ctx.Done()

	return nil
}

func ExtractPortrait(
	inputLocation string,
	width int,
	height int) (string, error) {

	log.Info("extracting portrait")

	outputFilePath := "/tmp/" + uuid.New().String() + ".jpg"

	cmd := exec.Command(
		"caire",
		"-in", inputLocation,
		"-out", outputFilePath,
		"-width="+strconv.Itoa(width),
		"-height="+strconv.Itoa(height),
		"-perc=1",
		"-square=0",
		"-scale=1",
		"-blur=0",
		"-sobel=0",
		"-debug=0",
		"-face=1",
		"-cc=./data/facefinder",
	)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	log.Info("extracted portrait")
	return outputFilePath, nil
}
