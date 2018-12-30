package main

import (
	"bytes"
	"flag"
	"fmt"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/crop/update_status"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"net/http"
	"os"
)

type Task struct {
	ID      string `jsonapi:"primary,crop_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}

func main() {
	initMinio()
	go startPrometheus()
	startFaktory()
}

func initMinio() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))
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
	fmt.Println("starting faktory")
	mgr := worker.NewManager()
	mgr.Use(func(perform worker.Handler) worker.Handler {
		return func(ctx worker.Context, job *faktory.Job) error {
			log.Printf("Starting work on job %s of type %s with custom %v\n", ctx.Jid(), ctx.JobType(), job.Custom)
			err := perform(ctx, job)
			log.Printf("Finished work on job %s with error %v\n", ctx.Jid(), err)
			return err
		}
	})
	mgr.Register("crop", convertTask)
	mgr.Queues = []string{"crop"}
	var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
	})
	mgr.Run()
}

func convertTask(ctx worker.Context, args ...interface{}) error {
	requests.With(prometheus.Labels{"controller": "crop", "status": "fetched"}).Inc()
	log.Info("Working on job %s\n", ctx.Jid())

	task := new(Task)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		log.Error("failed to deserialize task", args)
		return err
	}

	update_status.NotifyAboutProcessingStart(task.ID)

	err = handle(task)
	if err != nil {
		_ = ctx.Err()
	}

	update_status.NotifyAboutCompletion(task.ID)
	ctx.Done()
	requests.With(prometheus.Labels{"controller": "crop", "status": "completed"}).Inc()

	return nil
}

func handle(t *Task) error {
	downloadedFilePath, err := DownloadFile(t.ImageId)
	if err != nil {
		return err
	}
	croppedFilePath, err := CropImage(downloadedFilePath, t.Width, t.Height)
	if err != nil {
		return err
	}
	_, err = minioconnector.UploadFileWithName(croppedFilePath, t.ID)
	return err
}

func DownloadFile(objectName string) (string, error) {
	return minioconnector.DownloadFile(objectName)
}

func CropImage(inputImg string, width int, height int) (string, error) {
	log.Info("starting to crop api_image")

	outputFilePath := "/tmp/downloaded" + uuid.New().String() + ".jpg"

	f, _ := os.Open(inputImg)
	img, _, _ := image.Decode(f)
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, width, height)
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedImg := img.(SubImager).SubImage(topCrop)
	f, err := os.Create(outputFilePath)
	defer f.Close()
	if err != nil {
		return "", err
	}
	err = jpeg.Encode(f, croppedImg, nil)
	if err != nil {
		return "", err
	}

	log.Info("finished cropping api_image")
	return outputFilePath, nil
}
