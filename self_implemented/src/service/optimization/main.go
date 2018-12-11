package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/kmanuel/minioconnector"
	"image"
	"image/jpeg"
	"log"
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
	log.Println(http.ListenAndServe(":8080", router))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	downloadedFilePath := minioconnector.DownloadFile(task.In)

	outputFilePath := optimizeImage(downloadedFilePath, task.Width, task.Height)

	minioconnector.UploadFile(outputFilePath)

	log.Printf("finished")
}

func optimizeImage(inputFile string, width int, height int) string {
	outputFilePath := "/tmp/" + uuid.New().String() + ".jpg"

	f, _ := os.Open(inputFile)
	img, _, _ := image.Decode(f)
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, width, height)
	// The crop will have the requested aspect ratio, but you need to copy/scale it yourself
	fmt.Printf("Top crop: %+v\n", topCrop)
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

	return outputFilePath
}
