package main

import (
	"encoding/json"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
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
	godotenv.Load()
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))

	//router := mux.NewRouter()
	//router.HandleFunc("/", handleRequest).Methods("POST")
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
	mgr.Register("crop", convertTask)
	mgr.Queues = []string{"crop"}
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
		width, _ := strconv.Atoi(strings["width"].(string))
		height, _ := strconv.Atoi(strings["height"].(string))
		handle(strings["in"].(string), width, height)
	}

	return nil
}


func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Info("received new request")
	request := parseRequest(w, r)
	inputFileId := request.In
	width := request.Width
	height := request.Height
	handle(inputFileId, width, height)

	log.Info("finished request handling")
}

func handle(inputFileId string, width int, height int) {
	downloadedFilePath := DownloadFile(inputFileId)
	croppedFilePath := CropImage(downloadedFilePath, width, height)
	minioconnector.UploadFile(croppedFilePath)
}

func parseRequest(w http.ResponseWriter, r *http.Request) Request {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)
	return task
}
func DownloadFile(objectName string) string {
	return minioconnector.DownloadFile(objectName)
}

func CropImage(inputImg string, width int, height int) string {
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
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, croppedImg, nil)

	log.Info("finished cropping image")
	return outputFilePath
}
