package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kmanuel/minioconnector"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

type Request struct {
	In string `json:"in,omitempty"`
	Width int `json:"width"`
	Height int `json:"height"`
}

func main() {
	minioconnector.Init()

	router := mux.NewRouter()
	router.HandleFunc("/", handleRequest).Methods("POST")
	log.Println(http.ListenAndServe(":8082", router))
}

func handleRequest(w http.ResponseWriter, r * http.Request) {
	request := parseRequest(w, r)
	downloadedFilePath := DownloadFile(request.In)
	croppedFilePath := CropImage(downloadedFilePath, request.Width, request.Height)
	minioconnector.UploadFile(croppedFilePath)

	log.Printf("finished")
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

	return outputFilePath
}
