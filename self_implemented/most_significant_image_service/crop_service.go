package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

const InputImageLocation = "../test/input/"
const OutputImageLocation = "../test/output/"

type Request struct {
	In string `json:"in,omitempty"`
	Out string `json:"out,omitempty"`
	Width int `json:"width"`
	Height int `json:"height"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/crop", CropImage).Methods("POST")
	log.Println(http.ListenAndServe(":8082", router))
}

func CropImage(w http.ResponseWriter, r * http.Request) {
	json.NewEncoder(w).Encode("received request")
	var task Request
	_ = json.NewDecoder(r.Body).Decode(&task)

	f, _ := os.Open(InputImageLocation + task.In)
	img, _, _ := image.Decode(f)

	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, 250, 250)

	// The crop will have the requested aspect ratio, but you need to copy/scale it yourself
	fmt.Printf("Top crop: %+v\n", topCrop)

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedImg := img.(SubImager).SubImage(topCrop)

	f, err := os.Create(OutputImageLocation + task.Out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, croppedImg, nil)

	log.Printf("finished")
}
