package main

import (
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/portrait/update_status"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
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

	strings, ok := args[0].(map[string]interface{})
	if !ok {
		log.Error("couldnt convert args[0]")
	} else {
		update_status.NotifyAboutProcessingStart(strings["id"].(string))

		inputLocation := strings["in"].(string)

		downloadedFilePath := minioconnector.DownloadFile(inputLocation)

		width, _ := strconv.Atoi(strings["width"].(string))
		height, _ := strconv.Atoi(strings["height"].(string))

		outputFilePath := ExtractPortrait(downloadedFilePath, width, height)

		minioconnector.UploadFile(outputFilePath)

		update_status.NotifyAboutCompletion(strings["id"].(string))
	}

	return nil
}

func ExtractPortrait(
	inputLocation string,
	width int,
	height int) string {

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
	cmd.Run()

	log.Info("extracted portrait")
	return outputFilePath
}
