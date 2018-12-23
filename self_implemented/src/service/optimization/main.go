package main

import (
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/optimization/update_status"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"os"
	"strconv"
)

type Request struct {
	In     string `json:"in,omitempty"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func main() {
	godotenv.Load()
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))

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
		log.Error("couldnt convert args[0]")
	} else {
		update_status.NotifyAboutProcessingStart(strings["id"].(string))

		downloadedFilePath := minioconnector.DownloadFile(strings["in"].(string))

		width, _ := strconv.Atoi(strings["width"].(string))
		height, _ := strconv.Atoi(strings["height"].(string))
		outputFilePath := optimizeImage(downloadedFilePath, width, height)

		minioconnector.UploadFile(outputFilePath)

		update_status.NotifyAboutCompletion(strings["id"].(string))
	}

	return nil
}

func optimizeImage(inputFile string, width int, height int) string {
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
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, croppedImg, nil)

	log.Info("optimized image")
	return outputFilePath
}
