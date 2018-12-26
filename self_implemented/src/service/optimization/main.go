package main

import (
	"flag"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/optimization/update_status"
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
	mgr := worker.NewManager()
	mgr.Use(func(perform worker.Handler) worker.Handler {
		return func(ctx worker.Context, job *faktory.Job) error {
			log.Printf("Starting work on job %s of type %s with custom %v\n", ctx.Jid(), ctx.JobType(), job.Custom)
			err := perform(ctx, job)
			log.Printf("Finished work on job %s with error %v\n", ctx.Jid(), err)
			return err
		}
	})
	mgr.Register("optimization", convertTask)
	mgr.Queues = []string{"optimization"}
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
		_ = ctx.Err()
	} else {
		taskId := strings["id"].(string)
		update_status.NotifyAboutProcessingStart(taskId)

		downloadedFilePath, err := minioconnector.DownloadFile(strings["in"].(string))
		if err != nil {
			_ = ctx.Err()
			return nil
		}

		width, err := strconv.Atoi(strings["width"].(string))
		if err != nil {
			_ = ctx.Err()
			return nil
		}

		height, err := strconv.Atoi(strings["height"].(string))
		if err != nil {
			_ = ctx.Err()
			return nil
		}

		outputFilePath, err := optimizeImage(downloadedFilePath, width, height)
		if err != nil {
			_ = ctx.Err()
			return nil
		}

		_, err = minioconnector.UploadFileWithName(outputFilePath, taskId)
		if err != nil {
			_ = ctx.Err()
			return nil
		}

		update_status.NotifyAboutCompletion(taskId)

		ctx.Done()
	}

	return nil
}

func optimizeImage(inputFile string, width int, height int) (string, error) {
	log.Info("optimizing image")
	outputFilePath := "/tmp/" + uuid.New().String() + ".jpg"

	f, _ := os.Open(inputFile)
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


	log.Info("optimized image")
	return outputFilePath, nil
}
