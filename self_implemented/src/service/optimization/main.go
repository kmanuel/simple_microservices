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
	router.HandleFunc("/", HandleRequest).Methods("POST")
	log.Info(http.ListenAndServe(":8080", router))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	log.Info("received request")
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	downloadedFilePath := minioconnector.DownloadFile(task.In)

	outputFilePath := optimizeImage(downloadedFilePath, task.Width, task.Height)

	minioconnector.UploadFile(outputFilePath)

	log.Info("finished request")
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
