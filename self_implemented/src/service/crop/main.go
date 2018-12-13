package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"net/http"
	"os"
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

	router := mux.NewRouter()
	router.HandleFunc("/", handleRequest).Methods("POST")
	log.Info(http.ListenAndServe(":8080", router))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Info("received new request")
	request := parseRequest(w, r)
	downloadedFilePath := DownloadFile(request.In)
	croppedFilePath := CropImage(downloadedFilePath, request.Width, request.Height)
	minioconnector.UploadFile(croppedFilePath)

	log.Info("finished request handling")
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
