package main

import (
	"flag"
	"fmt"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
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
	"strconv"
)

type Request struct {
	In     string `json:"in,omitempty"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
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
		[]string{"service", "status"},
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
	log.Println("did register")
	mgr.Queues = []string{"crop"}
	var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
	})
	// Start processing jobs, this method does not return
	fmt.Println("running mgr")
	mgr.Run()
	fmt.Println("started faktory")
}

func convertTask(ctx worker.Context, args ...interface{}) error {
	requests.With(prometheus.Labels{"service":"crop", "status": "fetched"}).Inc()

	log.Info("Working on job %s\n", ctx.Jid())
	strings, ok := args[0].(map[string]interface{})
	if !ok {
		log.Error("couldnt convert args[0]")
		ctx.Err()
	} else {
		update_status.NotifyAboutProcessingStart(strings["id"].(string))

		width, _ := strconv.Atoi(strings["width"].(string))
		height, _ := strconv.Atoi(strings["height"].(string))
		err := handle(strings["id"].(string), strings["in"].(string), width, height)
		if err != nil {
			ctx.Err()
		}

		update_status.NotifyAboutCompletion(strings["id"].(string))
		ctx.Done()
		requests.With(prometheus.Labels{"service":"crop", "status": "completed"}).Inc()
	}

	return nil
}

func handle(taskId string, inputFileId string, width int, height int) error {
	downloadedFilePath := DownloadFile(inputFileId)
	croppedFilePath, err := CropImage(downloadedFilePath, width, height)
	if err != nil {
		return err
	}
	minioconnector.UploadFileWithName(croppedFilePath, taskId)
	return nil
}

func DownloadFile(objectName string) string {
	return minioconnector.DownloadFile(objectName)
}

func CropImage(inputImg string, width int, height int) (string, error) {
	log.Info("starting to crop image")

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
	if err != nil {
		return "", err
	}
	defer f.Close()
	err = jpeg.Encode(f, croppedImg, nil)
	if err != nil {
		return "", err
	}

	log.Info("finished cropping image")
	return outputFilePath, nil
}
