package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

const InputImgFolder = "../test/input/"
const OutputImgFolder = "../test/output/"

type Request struct {
	In string `json:"in,omitempty"`
	Out string `json:"out,omitempty"`
	Width int `json:"width"`
	Height int `json:"height"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handleRequest).Methods("POST")
	log.Println(http.ListenAndServe(":8082", router))
}

func handleRequest(w http.ResponseWriter, r * http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)
	inputImg := InputImgFolder + task.In
	outputImg := OutputImgFolder + task.Out

	CropImage(inputImg, outputImg)

	log.Printf("finished")
}

func CropImage(inputImg string, outputImg string) {
	f, _ := os.Open(inputImg)
	img, _, _ := image.Decode(f)
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, 250, 250)
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedImg := img.(SubImager).SubImage(topCrop)
	f, err := os.Create(outputImg)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, croppedImg, nil)
}
